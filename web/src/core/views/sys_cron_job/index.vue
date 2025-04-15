<template>
  <el-card class="container" shadow="never" v-loading="pageLoading">
    <div class="mb-[10px] ">
      <el-button v-if="multipleSelection.length <= 0" type="primary" icon="Refresh" :loading="pageLoading" @click="getPageData">刷新</el-button>
      <el-button v-if="multipleSelection.length <= 0" type="primary" icon="Plus" @click="()=>handleAdd({enable:true})">新增</el-button>
      <el-button v-if="multipleSelection.length > 0" type="danger" icon="Delete" @click="()=>handleDelete(multipleSelection)">批量删除</el-button>
    </div>
    <el-table :data="tableData"  @selection-change="handleSelectionChange">
      <el-table-column type="selection"></el-table-column>
      <el-table-column label="编号" prop="id" :width="100"></el-table-column>
      <el-table-column label="名称" prop="name"></el-table-column>
      <el-table-column label="执行任务" prop="handler_key">
        <template #default="{ row }">
          <el-tag>{{ handlers[row.handler_key].name }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Cron表达式" prop="cron"></el-table-column>
      <el-table-column label="最大并发" prop="max_concurrency">
        <template #default="{ row }">
          <el-tag type="primary">{{ row.max_concurrency }}</el-tag>
<!--          <el-tag v-if="row.allow_concurrency" type="success">是</el-tag>-->
<!--          <el-tag v-else type="danger">否</el-tag>-->
        </template>
      </el-table-column>
      <el-table-column label="是否启用" prop="enable">
        <template #default="{ row }">
          <el-tag v-if="row.enable" type="success">运行中</el-tag>
          <el-tag v-else type="danger">未运行</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" :width="160">
        <template #default="{ row }">
          <el-button-group class="table-btn-group">
            <el-button type="primary" icon="Edit" text @click="()=>handleAdd(row)">修改</el-button>
            <el-button type="danger" icon="Delete" text @click="()=>handleDelete([row.id])">删除</el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>
    <div class="mt-[15px] flex justify-center">
      <el-pagination
          v-model:current-page="queryForm.page"
          v-model:page-size="queryForm.limit"
          :page-sizes="[10, 30, 50, 100]"
          background
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @change="getPageData"
      />
    </div>
    <FormDialog v-model="dialogOpen" v-model:form="submitForm" :title="`${ submitForm.id ? '修改' : '新增' }任务管理`" :on-confirm="handleSubmit">
      <el-form-item label="名称" prop="name" :rules="[{required:true,message:'名称不能为空'}]">
        <el-input v-model="submitForm.name" placeholder="请输入名称"></el-input>
      </el-form-item>
      <el-form-item label="执行任务" prop="handler_key" :rules="[{required:true,message:'执行任务不能为空'}]">
        <el-select v-model="submitForm.handler_key">
          <el-option v-for="(handler,key) in handlers" :label="handler.name" :value="key"></el-option>
        </el-select>
      </el-form-item>
      <template v-if="handlers[submitForm.handler_key]?.params">
        <el-form-item v-for="item in handlers[submitForm.handler_key].params" :label="item.name" :prop="`params.${item.key}`" :rules="[{required:item.required,message:`${item.name}是必需的`}]">
          <el-input v-if="item.type === 'string'" v-model="submitForm.params[item.key]" :placeholder="`请输入${item.name}`"></el-input>
          <el-input-number v-if="item.type === 'int'" v-model="submitForm.params[item.key]" :placeholder="item.name"></el-input-number>
          <el-switch v-if="item.type === 'bool'" v-model="submitForm.params[item.key]" :active-value="true" :inactive-value="false"></el-switch>
        </el-form-item>
      </template>
      <el-form-item label="Cron表达式" prop="cron" :rules="[{required:true,message:'Cron表达式不能为空'}]">
        <el-input v-model="submitForm.cron" placeholder="请输入Cron表达式"></el-input>
      </el-form-item>
      <el-form-item label="最大并发" prop="max_concurrency">
        <el-input-number v-model="submitForm.max_concurrency" :min="1"></el-input-number>
      </el-form-item>
      <el-form-item label="是否启用" prop="enable">
        <el-switch v-model="submitForm.enable" :active-value="true" :inactive-value="false"></el-switch>
      </el-form-item>
    </FormDialog>
  </el-card>
</template>

<script setup lang="ts">
import { SysCronApi } from "../../apis/sys_cron";
import { ElMessage,ElMessageBox } from "element-plus";
import FormDialog from "../../../components/core/FormDialog.vue";
import {onMounted, ref} from "vue";

const dialogOpen = ref(false)
const queryForm = ref({
  page:1,
  limit:10
})
const pageLoading = ref(true)
const total = ref(0)
const multipleSelection:any = ref([])
const handlers:any = ref([])
const tableData = ref([])
const submitForm:any = ref({})

const getHandlers = async () => {
  const response = await SysCronApi.GetRegisteredHandler()
  handlers.value = response.data
}

const getPageData = async () => {
  pageLoading.value = true
  multipleSelection.value = []
  try {
    const response = await SysCronApi.List(queryForm.value)
    tableData.value = response.data.list
    total.value = response.data.total
  }catch (e) {
    console.log(e)
  }
  pageLoading.value = false
}


const handleSelectionChange = (val:Array<any>) => {
  multipleSelection.value = val.map(item=>item.id)
}

const handleAdd = (defaultForm:any)=>{
  console.log(defaultForm.params ? JSON.parse(defaultForm.params)  : {})
  submitForm.value = {
    ...defaultForm,
    params: defaultForm.params ? JSON.parse(defaultForm.params)  : {},
  }
  dialogOpen.value = true
}

const handleDelete = (ids:Array<any>) => {
  ElMessageBox.confirm("您确认要删除选中的数据吗？","删除提示",{
    type:"error",
    beforeClose:async (action, instance, done)=>{
      if (action === "confirm") {
        instance.confirmButtonLoading = true
        try {
          await SysCronApi.Delete(ids)
          ElMessage.success("删除成功")
          getPageData()
        }catch (e){
          console.log(e)
        }
        instance.confirmButtonLoading = false
      }
      done()
    }
  })
}

const handleSubmit = async ()=>{
  if (!submitForm.value.id){
    await SysCronApi.Create(submitForm.value)
    ElMessage.success("创建成功")
  }else{
    await SysCronApi.Edit(submitForm.value)
    ElMessage.success("修改成功")
  }
  getPageData()
}

onMounted(()=>{
  getHandlers()
  getPageData()
})
</script>


<style scoped>

</style>

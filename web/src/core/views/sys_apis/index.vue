<template>
  <el-card class="container" shadow="never" v-loading="pageLoading">
    <div class="mb-[10px] ">
      <el-button v-if="multipleSelection.length <= 0" type="primary" icon="Refresh" :loading="pageLoading" @click="getPageData">刷新</el-button>
      <el-button v-if="multipleSelection.length <= 0" type="primary" @click="handleSyncAPi">
        <template #icon>
          <iconify-icon icon="iconoir:cloud-desync"></iconify-icon>
        </template>
        <span>同步API</span>
      </el-button>
      <el-button v-if="multipleSelection.length > 0" type="danger" icon="Delete" @click="()=>handleDelete(multipleSelection)">批量删除</el-button>
    </div>
    <el-table :data="tableData"  @selection-change="handleSelectionChange">
      <el-table-column type="selection"></el-table-column>
      <el-table-column label="编号" prop="id" :width="100"></el-table-column>
      <el-table-column label="路径" prop="path"></el-table-column>
      <el-table-column label="方法" prop="method">
        <template #default="{ row }">
          <el-tag :type="MethodType[row.method] || 'warning'">{{ row.method }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="分组" prop="group"></el-table-column>
      <el-table-column label="描述" prop="description"></el-table-column>
      <el-table-column label="操作" :width="160">
        <template #default="{ row }">
          <el-button-group class="table-btn-group">
            <el-button type="primary" icon="Edit" text @click="()=>handleShowFormDialog(row)">修改</el-button>
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
    <FormDialog v-model="dialogOpen" v-model:form="submitForm" :title="`${ submitForm.id ? '修改' : '新增' }API管理`" :on-confirm="handleSubmit">
      <el-form-item label="分组" prop="group" :rules="[{required:true,message:'分组不能为空'}]">
        <el-select
            v-model="submitForm.group"
            filterable
            allow-create
            default-first-option
            :reserve-keyword="false"
            placeholder="选择或新增分组"
        >
          <el-option
              v-for="item in groups"
              :key="item"
              :label="item"
              :value="item"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="描述" prop="description" :rules="[{required:true,message:'描述不能为空'}]">
        <el-input v-model="submitForm.description" placeholder="请输入描述"></el-input>
      </el-form-item>
      <el-form-item label="方法" prop="method" :rules="[{required:true,message:'方法不能为空'}]">
        <el-select v-model="submitForm.method" placeholder="请选择方法">
          <el-option v-for="item in Object.keys(MethodType)" :label="item" :value="item"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="路径" prop="path" :rules="[{required:true,message:'路径不能为空'}]">
        <el-input v-model="submitForm.path" placeholder="请输入路径"></el-input>
      </el-form-item>
    </FormDialog>
    <Sync ref="syncRef" @on-add="getPageData"></Sync>
  </el-card>
</template>

<script setup lang="ts">
import { SysApisApi } from "../../apis/sys_apis";
import { ElMessage,ElMessageBox } from "element-plus";
import FormDialog from "../../../components/core/FormDialog.vue";
import {onMounted, ref} from "vue";
import Sync from "./sync_api.vue";
import {MethodType} from "./method_type.ts";
const syncRef = ref();
const dialogOpen = ref(false)
const queryForm = ref({
  page:1,
  limit:10
})
// 初始化form数据
const initForm = ()=>{
    return {
        group:"",
        description:"",
        method:"",
        path:"",
    }
}
const pageLoading = ref(true)
const total = ref(0)
const multipleSelection:any = ref([])
const tableData = ref([])
const submitForm:any = ref(initForm())
const groups:any = ref([])
const getPageData = async () => {
  pageLoading.value = true
  multipleSelection.value = []
  try {
    const response = await SysApisApi.List(queryForm.value)
    tableData.value = response.data.list
    total.value = response.data.total
    SysApisApi.GetGroups().then(res=>{
      groups.value = res.data
    })
  }catch (e) {
    console.log(e)
  }
  pageLoading.value = false
}


const handleSelectionChange = (val:Array<any>) => {
  multipleSelection.value = val.map(item=>item.id)
}

const handleShowFormDialog = (defaultForm:any)=>{
  submitForm.value = {
    ...defaultForm
  }
  dialogOpen.value = true
}

const handleSyncAPi = ()=>{
  syncRef.value.show()
}

const handleDelete = (ids:Array<any>) => {
  ElMessageBox.confirm("您确认要删除选中的数据吗？","删除提示",{
    type:"error",
    beforeClose:async (action, instance, done)=>{
      if (action === "confirm") {
        instance.confirmButtonLoading = true
        try {
          await SysApisApi.Delete(ids)
          ElMessage.success("删除成功")
          getPageData()
        }catch (e){
          console.log(e)
        }
        done()
        instance.confirmButtonLoading = false
      } else if (!instance.confirmButtonLoading){
        done()
      }
    }
  })
}

const handleSubmit = async ()=>{
  if (!submitForm.value.id){
    await SysApisApi.Create(submitForm.value)
    ElMessage.success("创建成功")
  }else{
    await SysApisApi.Edit(submitForm.value)
    ElMessage.success("修改成功")
  }
  getPageData()
}

onMounted(()=>{
  getPageData()
})
</script>


<style scoped>

</style>

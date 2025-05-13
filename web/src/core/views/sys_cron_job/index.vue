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
      <el-table-column label="执行周期" prop="cron">
        <template #default="{ row }">
          <span class="text-[14px]">{{ transSpecToStr(row.cron) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="是否启用" prop="enable">
        <template #default="{ row }">
          <el-tag v-if="row.enable" type="success">运行中</el-tag>
          <el-tag v-else type="danger">未运行</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="最近执行时间" prop="last_run_time">
        <template #default="{ row }">
          <span class="text-[12px]">{{row.last_run_time ? formatTime(row.last_run_time) : '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="下次执行时间" prop="next_run_time">
        <template #default="{ row }">
          <span class="text-[12px]">{{ row.next_run_time ? formatTime(row.next_run_time) : '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" :width="260">
        <template #default="{ row }">
          <el-button-group class="table-btn-group">
            <el-button type="primary" icon="Edit" text @click="()=>handleAdd(row)">编辑</el-button>
            <el-button type="primary" icon="Tickets" text @click="()=>handleShowLog(row)">查看日志</el-button>
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
    <FormDialog v-model="dialogOpen" v-model:form="submitForm" :title="`${ submitForm.id ? '修改' : '新增' }计划任务`" :on-confirm="handleSubmit">
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
      <el-form-item label="执行周期" prop="cron" :rules="[{validator:checkCron}]">
        <el-row :gutter="10" class="w-full">
          <el-col :span="6">
            <el-select v-model="submitForm.specType" class="w-full">
              <el-option v-for="(item,key) in cronMethod" :label="item.label" :value="key"></el-option>
            </el-select>
          </el-col>
          <el-col v-if="needShow(submitForm.specType,'week')" :span="6">
            <el-select v-model="submitForm.week">
              <el-option v-for="(value,index) in weeks" :label="value" :value="index + 1"></el-option>
            </el-select>
          </el-col>
          <el-col v-if="needShow(submitForm.specType,'day')" :span="6">
            <el-input-number v-model="submitForm.day" :controls="false" class="w-[100%!important]" :min="1" :max="31" :value-on-clear="1">
              <template #suffix>日</template>
            </el-input-number>
          </el-col>
          <el-col v-if="needShow(submitForm.specType,'hour')" :span="6">
            <el-input-number v-model="submitForm.hour" :controls="false" class="w-[100%!important]" :min="0" :max="24" :value-on-clear="0">
              <template #suffix>小时</template>
            </el-input-number>
          </el-col>
          <el-col v-if="needShow(submitForm.specType,'minute')" :span="6">
            <el-input-number v-model="submitForm.minute" :controls="false" class="w-[100%!important]" :min="submitForm.specType =='perNMinute' ? 1 : 0 " :max="59" :value-on-clear="submitForm.specType =='perNMinute' ? 1 : 0">
              <template #suffix>分钟</template>
            </el-input-number>
          </el-col>
          <el-col v-if="needShow(submitForm.specType,'second')" :span="6">
            <el-input-number v-model="submitForm.second" :controls="false" class="w-[100%!important]" :min="1" :max="3600" :value-on-clear="1">
              <template #suffix>秒</template>
            </el-input-number>
          </el-col>
        </el-row>
      </el-form-item>
      <el-form-item label="是否启用" prop="enable">
        <el-switch v-model="submitForm.enable" :active-value="true" :inactive-value="false"></el-switch>
      </el-form-item>
    </FormDialog>
    <Logs ref="logRef"></Logs>
  </el-card>
</template>

<script setup lang="ts">
import { SysCronJobApi } from "../../apis/sys_cron_job";
import { ElMessage,ElMessageBox } from "element-plus";
import FormDialog from "../../../components/core/FormDialog.vue";
import {onMounted, ref} from "vue";
import {cronMethod, needShow, transObjToSpec, transSpecToObj, transSpecToStr, weeks} from "./helper.ts";
import {formatTime} from "../../../utils/formatTime.ts";
import Logs from "./logs.vue";
const logRef = ref();
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

const checkCron = (_rule:any, _value:any, callback: any) => {
  let result = transObjToSpec(
      submitForm.value.specType,
      submitForm.value.week,
      submitForm.value.day,
      submitForm.value.hour,
      submitForm.value.minute,
      submitForm.value.second,
  )
  if (!result){
    callback("执行周期设置不正确，请检查")
  }else{
    callback()
  }
}


const getHandlers = async () => {
  const response = await SysCronJobApi.GetRegisteredHandler()
  handlers.value = response.data
}

const getPageData = async () => {
  pageLoading.value = true
  multipleSelection.value = []
  try {
    const response = await SysCronJobApi.List(queryForm.value)
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
  submitForm.value = {
    ...defaultForm,
    params: defaultForm.params ? JSON.parse(defaultForm.params)  : {},
    ...transSpecToObj(defaultForm.cron),
    max_concurrency:defaultForm.max_concurrency || 1
  }
  dialogOpen.value = true
}

const handleShowLog = (row:any)=>{
  // logsOpen.value = true
  logRef.value.showLog(row)
}

const handleDelete = (ids:Array<any>) => {
  ElMessageBox.confirm("您确认要删除选中的数据吗？","删除提示",{
    type:"error",
    beforeClose:async (action, instance, done)=>{
      if (action === "confirm") {
        instance.confirmButtonLoading = true
        try {
          await SysCronJobApi.Delete(ids)
          ElMessage.success("删除成功")
          getPageData()
        }catch (e){
          console.log(e)
        }
        done()
        instance.confirmButtonLoading = false
      }else{
        done()
      }
    }
  })
}

const handleSubmit = async ()=>{
  submitForm.value.cron = transObjToSpec(
      submitForm.value.specType,
      submitForm.value.week,
      submitForm.value.day,
      submitForm.value.hour,
      submitForm.value.minute,
      submitForm.value.second,
  )
  if (!submitForm.value.id){
    await SysCronJobApi.Create(submitForm.value)
    ElMessage.success("创建成功")
  }else{
    await SysCronJobApi.Edit(submitForm.value)
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

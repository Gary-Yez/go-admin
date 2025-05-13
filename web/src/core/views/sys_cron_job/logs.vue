<template>
  <el-drawer v-model="show" size="50vw" title="执行日志">
    <div class="mb-[10px] ">
      <el-button type="primary" icon="Refresh" :loading="pageLoading" @click="getPageData">刷新</el-button>
    </div>
    <el-table :data="tableData" v-loading="pageLoading" row-key="id">
      <el-table-column label="ID" prop="id" :width="100">
        <template #default="{ row }">
          <span class="text-[12px]">{{ row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column label="开始时间" prop="start_time" :width="150">
        <template #default="{ row }">
          <span class="text-[12px]">{{ formatTime(row.start_time) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="结束时间" prop="end_time" :width="150">
        <template #default="{ row }">
          <span class="text-[12px]">{{ formatTime(row.end_time) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="用时" :width="100">
        <template #default="{ row }">
          <el-tag type="primary">{{ ((+new Date(row.end_time) - +new Date(row.start_time))/ 1000).toFixed(2)}} 秒</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="执行结果" :width="100">
        <template #default="{ row }">
          <el-tag v-if="!row.error" type="success">成功</el-tag>
          <el-tag v-else type="danger">失败</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="失败原因">
        <template #default="{ row }">
          <span class="text-[12px]">{{ row.error || "success" }}</span>
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
  </el-drawer>
</template>

<script setup lang="ts">
import {onMounted, ref} from "vue";
  import {SysCronJobApi} from "../../apis/sys_cron_job.ts";
import {formatTime} from "../../../utils/formatTime.ts";
  const props:any = defineProps({
    job:{
      type:Object,
    }
  })
  const pageLoading = ref(true);
  const show = ref(false);
  const tableData = ref([])
  const total = ref(0)
  const queryForm:any = ref({
    page:1,
    limit:10,
  })

  onMounted(()=>{
    console.log(props.job)
  })

  const getPageData = async () => {
    pageLoading.value = true
    try {
      const response = await SysCronJobApi.GetLogs(queryForm.value)
      tableData.value = response.data.list
      total.value = response.data.total
    }catch(error) {
      console.log(error)
    }
    pageLoading.value = false
  }

  defineExpose({
    showLog(row:any){
      tableData.value = []
      total.value = 0
      queryForm.value = {
        page: 1,
        limit: 10,
        filters:[{
          field:"job_id",
          operator:"=",
          value:`${row.id}`
        }]
      }
      show.value = true
      getPageData()
    },
  })
</script>

<style scoped lang="less">
</style>

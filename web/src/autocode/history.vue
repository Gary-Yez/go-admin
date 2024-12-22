<template>
  <el-card class="container" shadow="never">
    <div class="mb-[10px] ">
      <el-button type="primary" icon="Refresh" :loading="pageLoading" @click="getPageData">刷新</el-button>
      <el-button type="primary" icon="Plus" @click="handleGoAdd">生成代码</el-button>
    </div>
    <el-table :data="tableData">
      <el-table-column label="ID" prop="id" :width="100"></el-table-column>
      <el-table-column label="模块名称" prop="module_name"></el-table-column>
      <el-table-column label="结构体名称" prop="model_name"></el-table-column>
      <el-table-column label="创建时间" prop="created_at" :width="200">
        <template #default="{ row }">
          <span>{{ formatTime(row.created_at) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="更新时间" prop="created_at" :width="200">
        <template #default="{ row }">
          <span>{{ formatTime(row.created_at) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" :width="160">
        <template #default="{ row }">
          <el-button-group class="table-btn-group">
            <el-button type="primary" icon="Edit" text @click="()=>handleEdit(row)">修改</el-button>
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
  </el-card>
</template>

<script setup lang="ts">
  import {onMounted, ref} from "vue";
  import {SysAutocodeApi} from "../apis/sys_autocode.ts";
  import {useRouter} from "vue-router";
  import {formatTime} from "../utils/formatTime.ts";
  import { ElMessage,ElMessageBox } from "element-plus"
  const router = useRouter()
  const pageLoading = ref(true)
  const queryForm = ref({
    page:1,
    limit:10
  })
  const tableData = ref([])
  const total = ref(0)

  const getPageData = async () => {
    pageLoading.value = true
    try {
      const response = await SysAutocodeApi.History(queryForm.value)
      tableData.value = response.data.list
      total.value = response.data.total
    }catch(err) {
      console.log(err)
    }
    pageLoading.value = false
  }

  const handleGoAdd = () => {
    router.push("/dashboard/sys_autocode")
  }

  onMounted(()=>{
    getPageData()
  })

  const handleEdit = (row: any) => {
    router.push({
      path:"/dashboard/sys_autocode",
      query:{id:row.id}
    })
  }

  const handleDelete = (ids:Array<any>) => {
    ElMessageBox.confirm("删除生成记录并不会删除已生成的文件，您确认要删除该记录吗？","删除提示",{
      type:"error",
      beforeClose:async (action:any, instance:any, done:any)=>{
        if (action === "confirm") {
          instance.confirmButtonLoading = true
          try {
            await SysAutocode.Delete(ids)
            ElMessage.success("删除成功")
            getPageData().then()
          }catch (e){
            console.log(e)
          }
          instance.confirmButtonLoading = false
        }
        done()
      }
    })
  }
</script>

<style scoped>

</style>

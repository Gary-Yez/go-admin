<template>
  <el-tabs v-model="activeName" type="border-card">
<!--    <el-tab-pane label="站点设置" name="site">-->
<!--      <el-form label-position="top">-->
<!--        <el-form-item v-for="item in tableData" :label="item.name" :prop="item.key" :key="item.id">-->
<!--          <el-input v-if="item.type === 'string'"></el-input>-->
<!--          <el-input-number v-else-if="item.type === 'int'"></el-input-number>-->
<!--          <el-switch v-else-if="item.type === 'bool'"></el-switch>-->
<!--        </el-form-item>-->
<!--      </el-form>-->
<!--    </el-tab-pane>-->
  </el-tabs>
<!--  <el-card class="container" shadow="never" v-loading="pageLoading">-->
<!--    <div class="mb-[15px]">-->
<!--      <el-button type="primary" @click="isManage = !isManage">切换{{ isManage ? '视图' : '管理'  }}模式</el-button>-->
<!--    </div>-->
<!--    <template v-if="isManage">-->
<!--      <div class="mb-[10px] ">-->
<!--        <el-button v-if="multipleSelection.length <= 0" type="primary" icon="Refresh" :loading="pageLoading" @click="getPageData">刷新</el-button>-->
<!--        <el-button v-if="multipleSelection.length <= 0" type="primary" icon="Plus" @click="()=>handleAdd({})">新增</el-button>-->
<!--        <el-button v-if="multipleSelection.length > 0" type="danger" icon="Delete" @click="()=>handleDelete(multipleSelection)">批量删除</el-button>-->
<!--      </div>-->
<!--      <el-table :data="tableData"  @selection-change="handleSelectionChange">-->
<!--        <el-table-column type="selection"></el-table-column>-->
<!--        <el-table-column label="编号" prop="id" :width="100"></el-table-column>-->
<!--        <el-table-column label="变量键" prop="key"></el-table-column>-->
<!--        <el-table-column label="变量名称" prop="name"></el-table-column>-->
<!--        <el-table-column label="变量类型" prop="type">-->
<!--          <template #default="{ row }">-->
<!--            <el-tag type="primary">{{ varType[row.type] || '未知' }}</el-tag>-->
<!--          </template>-->
<!--        </el-table-column>-->
<!--        <el-table-column label="变量值" prop="value"></el-table-column>-->
<!--        <el-table-column label="备注" prop="remark"></el-table-column>-->
<!--        <el-table-column label="操作" :width="160">-->
<!--          <template #default="{ row }">-->
<!--            <el-button-group class="table-btn-group">-->
<!--              <el-button type="primary" icon="Edit" text @click="()=>handleAdd(row)">修改</el-button>-->
<!--              <el-button type="danger" icon="Delete" text @click="()=>handleDelete([row.id])">删除</el-button>-->
<!--            </el-button-group>-->
<!--          </template>-->
<!--        </el-table-column>-->
<!--      </el-table>-->
<!--      <div class="mt-[15px] flex justify-center">-->
<!--        <el-pagination-->
<!--            v-model:current-page="queryForm.page"-->
<!--            v-model:page-size="queryForm.limit"-->
<!--            :page-sizes="[10, 30, 50, 100]"-->
<!--            background-->
<!--            layout="total, sizes, prev, pager, next, jumper"-->
<!--            :total="total"-->
<!--            @change="getPageData"-->
<!--        />-->
<!--      </div>-->
<!--      <FormDialog v-model="dialogOpen" v-model:form="submitForm" :title="`${ submitForm.id ? '修改' : '新增' }全局变量`" :on-confirm="handleSubmit">-->
<!--        <el-form-item label="变量键" prop="key" :rules="[{required:true,message:'变量键不能为空'}]">-->
<!--          <el-input v-model.trim="submitForm.key" :disabled="submitForm.id" placeholder="请输入变量键"></el-input>-->
<!--        </el-form-item>-->
<!--        <el-form-item label="变量名称" prop="name" :rules="[{required:true,message:'变量名称不能为空'}]">-->
<!--          <el-input v-model="submitForm.name" placeholder="请输入变量名称"></el-input>-->
<!--        </el-form-item>-->
<!--        <el-form-item label="变量类型" prop="type" :rules="[{required:true,message:'变量类型不能为空'}]">-->
<!--          <el-select v-model="submitForm.type" placeholder="请选择变量类型">-->
<!--            <el-option v-for="(label,value) in varType" :label="label" :value="value"></el-option>-->
<!--          </el-select>-->
<!--        </el-form-item>-->
<!--        <el-form-item label="变量值" prop="value">-->
<!--          <el-input v-model="submitForm.value" placeholder="请输入变量值"></el-input>-->
<!--        </el-form-item>-->
<!--        <el-form-item label="备注" prop="remark">-->
<!--          <el-input v-model="submitForm.remark" placeholder="请输入备注"></el-input>-->
<!--        </el-form-item>-->
<!--      </FormDialog>-->
<!--    </template>-->
<!--    <template v-else>-->
<!--      <el-form label-position="top">-->
<!--        <el-form-item v-for="item in tableData" :label="item.name" :prop="item.key" :key="item.id">-->
<!--          <el-input v-if="item.type === 'string'"></el-input>-->
<!--          <el-input-number v-else-if="item.type === 'int'"></el-input-number>-->
<!--          <el-switch v-else-if="item.type === 'bool'"></el-switch>-->
<!--        </el-form-item>-->
<!--      </el-form>-->
<!--    </template>-->
<!--  </el-card>-->
</template>

<script setup lang="ts">
import { SysGlobalVariableApi } from "../../apis/sys_global_variable.ts";
// import { ElMessage,ElMessageBox } from "element-plus";
// import FormDialog from "../../../components/core/FormDialog.vue";
import {onMounted, ref} from "vue";
const activeName = ref("site")
// const varType:any = ref({
//   "int":"整数型",
//   "string":"字符串",
//   "bool":"布尔值",
// });
// const isManage = ref(false);
// const dialogOpen = ref(false)
const queryForm = ref({
  page:1,
  limit:10
})
const pageLoading = ref(true)
const total = ref(0)
const multipleSelection:any = ref([])
const tableData:any = ref([])
// const submitForm:any = ref({})

const getPageData = async () => {
  pageLoading.value = true
  multipleSelection.value = []
  try {
    const response = await SysGlobalVariableApi.List(queryForm.value)
    tableData.value = response.data.list
    total.value = response.data.total
  }catch (e) {
    console.log(e)
  }
  pageLoading.value = false
}

//
// const handleSelectionChange = (val:Array<any>) => {
//   multipleSelection.value = val.map(item=>item.id)
// }
//
// const handleAdd = (defaultForm:any)=>{
//   submitForm.value = {
//     ...defaultForm
//   }
//   dialogOpen.value = true
// }
//
// const handleDelete = (ids:Array<any>) => {
//   ElMessageBox.confirm("您确认要删除选中的数据吗？","删除提示",{
//     type:"error",
//     beforeClose:async (action, instance, done)=>{
//       if (action === "confirm") {
//         instance.confirmButtonLoading = true
//         try {
//           await SysGlobalVariableApi.Delete(ids)
//           ElMessage.success("删除成功")
//           getPageData()
//         }catch (e){
//           console.log(e)
//         }
//         instance.confirmButtonLoading = false
//       }
//       done()
//     }
//   })
// }
//
// const handleSubmit = async ()=>{
//   if (!submitForm.value.id){
//     await SysGlobalVariableApi.Create(submitForm.value)
//     ElMessage.success("创建成功")
//   }else{
//     await SysGlobalVariableApi.Edit(submitForm.value)
//     ElMessage.success("修改成功")
//   }
//   getPageData()
// }

onMounted(()=>{
  getPageData()
})
</script>


<style scoped lang="less">

</style>

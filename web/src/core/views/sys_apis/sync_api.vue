<template>
  <FormDialog v-model="syncShow" title="同步API" :max-width="900" size="default" confirm-btn-text="确认同步" :on-confirm="handleConfirm">
    <template #default>
      <el-card header="新增API" shadow="never" >
        <el-table :data="pageData.newApis" v-loading="pageLoading">
          <el-table-column label="API路径" prop="path"></el-table-column>
          <el-table-column label="分组" prop="group">
            <template #default="{ row }">
              <el-select v-model="row.group" allow-create filterable default-first-option>
                <el-option v-for="item in groupList" :label="item" :value="item"></el-option>
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="描述" prop="description">
            <template #default="{ row }">
              <el-input v-model="row.description" placeholder="API描述"></el-input>
            </template>
          </el-table-column>
          <el-table-column label="请求" prop="method" :width="100">
            <template #default="{ row }">
              <el-tag :type="MethodType[row.method] || 'warning'">{{ row.method }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" :width="110">
            <template #default="{ row }">
              <el-button type="primary" link @click="handleAdd(row)">新增</el-button>
              <el-button type="primary" link @click="handleUpdateIgnore(row,true)">忽略</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
      <el-card class="mt-[20px]" header="待删除的API" shadow="never">
        <el-table :data="pageData.deleteApis" v-loading="pageLoading">
          <el-table-column label="API路径" prop="path"></el-table-column>
          <el-table-column label="分组" prop="group"></el-table-column>
          <el-table-column label="描述" prop="description"></el-table-column>
          <el-table-column label="请求" prop="method" :width="100">
            <template #default="{ row }">
              <el-tag :type="MethodType[row.method] || 'warning'">{{ row.method }}</el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
      <el-card class="mt-[20px]" header="忽略的API" shadow="never">
        <el-table :data="pageData.ignoreApis" v-loading="pageLoading">
          <el-table-column label="API路径" prop="path"></el-table-column>
          <el-table-column label="请求" prop="method">
            <template #default="{ row }">
              <el-tag :type="MethodType[row.method] || 'warning'">{{ row.method }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" :width="100">
            <template #default="{ row }">
              <el-button type="primary" link @click="handleUpdateIgnore(row,false)">取消忽略</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </template>
  </FormDialog>
</template>

<script setup lang="ts">
  import { ref } from "vue";
  import {SysApisApi} from "../../apis/sys_apis.ts";
  import {ElMessage} from "element-plus";
  import {MethodType} from "./method_type.ts";
  import FormDialog from "../../../components/core/FormDialog.vue";
  const emits = defineEmits(["onAdd"])
  const syncShow = ref(false);
  const pageLoading = ref(true)
  const pageData:any = ref({})
  const groupList = ref([])

  const getGroupList= async ()=>{
    const groups = await SysApisApi.GetGroups()
    groupList.value = groups.data
  }
  const getPageData = async () => {
    pageLoading.value = true
    const apis = await SysApisApi.SyncApi()
    pageData.value = apis.data
    pageLoading.value = false
  }

  const handleAdd = async (row:any)=>{
    if (!row.path || !row.description) return ElMessage.error("API分组和描述不能为空")
    await SysApisApi.Create(row)
    pageData.value.newApis = pageData.value.newApis.filter((item:any)=>{
      return !(item.path === row.path && item.method === row.method)
    })
    emits("onAdd")
    await getGroupList()
  }

  const handleUpdateIgnore = async (row:any,ignore:boolean)=>{
    await SysApisApi.UpdateIgnore(row.path,row.method,ignore)
    if (ignore){
      pageData.value.newApis = pageData.value.newApis.filter((item:any)=>{
        return !(item.path === row.path && item.method === row.method)
      })
      pageData.value.ignoreApis.push(row)
    }else {
      pageData.value.ignoreApis = pageData.value.ignoreApis.filter((item:any)=>{
        return !(item.path === row.path && item.method === row.method)
      })
      pageData.value.newApis.push(row)
    }
  }

  defineExpose({
    show(){
      getGroupList()
      getPageData()
      syncShow.value = true
    }
  })

  const handleConfirm = async () => {
    let pass = pageData.value.newApis.every((item:any)=>item.group && item.description)
    if (!pass){
      ElMessage.error("存在未填写分组和描述的API")
      throw new Error()
    }
    let createPromiseList = pageData.value.newApis.map((item:any)=>{
      return SysApisApi.Create({
       ...item,
      })
    })
    let deleteIds = pageData.value.deleteApis.map((item:any)=>item.id)
    if (createPromiseList.length > 0){
      await Promise.all(createPromiseList)
    }
    if (deleteIds.length > 0){
      await SysApisApi.Delete(deleteIds)
    }
    ElMessage.success("同步成功")
    getPageData()
  }
</script>

<style scoped>

</style>

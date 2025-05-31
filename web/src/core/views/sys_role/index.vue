<template>
  <el-card class="container" shadow="never">
    <div class="mb-[10px] ">
      <el-button type="primary" icon="Refresh" :loading="pageLoading" @click="getPageData">刷新</el-button>
      <el-button type="primary" icon="Plus" @click="()=>handleAdd({})">新增角色</el-button>
    </div>
    <el-table :data="tableData">
      <el-table-column label="编号" prop="id" :width="100"></el-table-column>
      <el-table-column label="名称" prop="name"></el-table-column>
      <el-table-column label="类型" prop="default">
        <template #default="{ row }">
          <el-tag v-if="row.default" type="primary">系统管理员</el-tag>
          <el-tag v-else type="success">自定义</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" :width="260">
        <template #default="{ row }">
          <el-button-group class="table-btn-group">
            <el-button type="primary" icon="Setting" :disabled="row.default" text @click="()=>handleEditPermission(row.id)">权限管理</el-button>
            <el-button type="primary" icon="Edit" text :disabled="row.default" @click="()=>handleAdd(row)">修改</el-button>
            <el-button type="danger" icon="Delete" text :disabled="row.default" @click="()=>handleDelete([row.id])">删除</el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>
    <FormDialog v-model="dialogOpen" v-model:form="submitForm" :title="submitForm.id ? '修改角色' : '新增角色'" :on-confirm="handleSubmit">
      <el-form-item label="角色名称" prop="name" :rules="[{required:true,message:'请输入角色名称'}]">
        <el-input v-model="submitForm.name" placeholder="请输入角色名称"></el-input>
      </el-form-item>
    </FormDialog>
    <FormDialog v-model="permissionOpen" v-model:form="permissionForm" title="权限管理" :on-confirm="handleSubmitPermission">
      <el-tabs type="border-card">
        <el-tab-pane label="菜单权限">
          <el-tree ref="menuTreeRef" :data="menus" show-checkbox default-expand-all node-key="id" :props="{label:'name'}"></el-tree>
        </el-tab-pane>
        <el-tab-pane label="API权限">
          <el-tree ref="apiTreeRef" :data="apis" show-checkbox default-expand-all node-key="id" :props="{label:'description'}">
            <template #default="{ data }">
              <div class="flex justify-between w-full">
                <span>{{ data.description }}</span>
                <span>{{ data.path }}</span>
              </div>
            </template>
          </el-tree>
        </el-tab-pane>
      </el-tabs>
    </FormDialog>
  </el-card>
</template>

<script setup lang="ts">
import {nextTick, onMounted, ref} from "vue";
import {SysRoleApi} from "../../apis/sys_role.ts";
import {SysMenuApi} from "../../apis/sys_menu.ts";
import { ElMessage,ElMessageBox } from "element-plus";
import FormDialog from "../../../components/core/FormDialog.vue";
import {SysApisApi} from "../../apis/sys_apis.ts";
const menuTreeRef = ref()
const apiTreeRef = ref()
const pageLoading = ref(true)
const menus = ref([])
const apis = ref([])
const apisIdMap:any = ref({})
const tableData = ref([])
const dialogOpen = ref(false)
const permissionOpen = ref(false)
const submitForm:any = ref({})
const permissionForm:any = ref({})

const getMenus = async ()=>{
  const response = await SysMenuApi.List()
  menus.value = response.data.list
}

const getApi = async ()=>{
  let apisMap:any = {}
  const response = await SysApisApi.List()
  response.data.list.forEach((item:any)=>{
    item.id = item.method+ item.path
    if (!apisMap[item.group]){
      apisMap[item.group] = {
        description: item.group,
        children: [item]
      }
    }else{
      apisMap[item.group].children.push(item)
    }
    apisIdMap.value[item.id] = item
  })
  apis.value = Object.values(apisMap)
}

const getPageData = async () => {
  pageLoading.value = true
  const response = await SysRoleApi.List()
  tableData.value = response.data.list
  pageLoading.value = false
}

const handleEditPermission = async (id:number)=>{
  await Promise.all([getApi(), getMenus()])
  const response = await SysRoleApi.Get(id)
  permissionForm.value = response.data
  permissionOpen.value = true
  await nextTick(()=>{
    apiTreeRef.value.setCheckedKeys(permissionForm.value.apis.map((item:any)=>item.method + item.path))
    menuTreeRef.value.setCheckedKeys(permissionForm.value.menus.map((item:any)=>item.id))
  })
}

const handleAdd = (defaultForm:any)=>{
  submitForm.value = {
    ...defaultForm,
  }
  dialogOpen.value = true
}

onMounted(()=>{
  getPageData()
})

const handleSubmitPermission = async () => {
  let menus = menuTreeRef.value.getCheckedKeys()
  let apis = apiTreeRef.value.getCheckedKeys().filter((item:any)=>item).map((item:any)=> {
    return {
      path:apisIdMap.value[item].path,
      method:apisIdMap.value[item].method,
    }
  })
  await SysRoleApi.UpdatePermission(permissionForm.value.id,menus,apis)
  ElMessage.success("权限修改成功")
}

const handleSubmit = async () => {
  if (!submitForm.value.id){
    await SysRoleApi.Create({
      ...submitForm.value,
    })
    ElMessage.success("创建成功")
  }else{
    await SysRoleApi.Edit(submitForm.value)
    ElMessage.success("修改成功")
  }
  getPageData()
}

const handleDelete = (ids:Array<any>) => {
  ElMessageBox.confirm("您确认要删除该角色吗？","删除提示",{
    type:"error",
    beforeClose:async (action, instance, done)=>{
      if (action === "confirm") {
        instance.confirmButtonLoading = true
        try {
          await SysRoleApi.Delete(ids)
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

</script>

<style scoped>

</style>

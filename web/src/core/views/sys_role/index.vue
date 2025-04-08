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
            <el-button type="primary" icon="Setting" :disabled="row.default" text @click="()=>handleAdd({id:row.id,editMenu:row.menus})">权限管理</el-button>
            <el-button type="primary" icon="Edit" text :disabled="row.default" @click="()=>handleAdd(row)">修改</el-button>
            <el-button type="danger" icon="Delete" text :disabled="row.default" @click="()=>handleDelete([row.id])">删除</el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>
    <FormDialog v-model="dialogOpen" v-model:form="submitForm" :title="submitForm.id ? '修改角色' : '新增角色'" :on-confirm="handleSubmit">
      <el-form-item v-if="!submitForm.editMenu" label="角色名称" prop="name" :rules="[{required:true,message:'请输入角色名称'}]">
        <el-input v-model="submitForm.name" placeholder="请输入角色名称"></el-input>
      </el-form-item>
      <el-form-item v-else label="角色权限" prop="menus">
        <el-tree ref="treeRef" :data="menus" show-checkbox default-expand-all node-key="id" :props="{label:'name'}">
        </el-tree>
      </el-form-item>
    </FormDialog>
  </el-card>
</template>

<script setup lang="ts">
import {nextTick, onMounted, ref} from "vue";
import {SysRoleApi} from "../../apis/sys_role.ts";
import {SysMenuApi} from "../../apis/sys_menu.ts";
import { ElMessage,ElMessageBox } from "element-plus";
import FormDialog from "../../../components/core/FormDialog.vue";
const treeRef = ref()
const pageLoading = ref(true)
const menus = ref([])
const tableData = ref([])
const dialogOpen = ref(false)
const submitForm:any = ref({})

const getMenus = async ()=>{
  const response = await SysMenuApi.List()
  menus.value = response.data.list
}

const getPageData = async () => {
  pageLoading.value = true
  const response = await SysRoleApi.List()
  tableData.value = response.data.list
  pageLoading.value = false
}

const handleAdd = (defaultForm:any)=>{
  submitForm.value = {
    ...defaultForm,
  }
  dialogOpen.value = true
  if (defaultForm.editMenu){
    nextTick(()=>{
      treeRef.value.setCheckedKeys(defaultForm.editMenu.map((item:any)=>item.id))
    })
  }
}

onMounted(()=>{
  getMenus()
  getPageData()
})

const handleSubmit = async () => {
  if (submitForm.value.editMenu){
    let menus = treeRef.value.getCheckedKeys()
    await SysRoleApi.UpdatePermission(submitForm.value.id,menus)
    ElMessage.success("权限修改成功")
  }else if (!submitForm.value.id){
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
        instance.confirmButtonLoading = false
      }
      done()
    }
  })
}

</script>

<style scoped>

</style>

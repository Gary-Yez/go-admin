<template>
  <el-card class="container" shadow="never">
    <div class="mb-[10px] ">
      <el-button type="primary" icon="Refresh" :loading="pageLoading" @click="getPageData">刷新</el-button>
      <el-button type="primary" icon="Plus" @click="()=>handleAdd({})">新增角色</el-button>
    </div>
    <el-table :data="tableData">
      <el-table-column label="编号" prop="id" :width="100"></el-table-column>
      <el-table-column label="名称" prop="name"></el-table-column>
      <el-table-column label="操作" :width="260">
        <template #default="{ row }">
          <el-button-group class="table-btn-group">
            <el-button type="primary" icon="Setting" text @click="()=>handleAdd({id:row.id,editMenu:row.menus})">权限管理</el-button>
            <el-button type="primary" icon="Edit" text @click="()=>handleAdd(row)">修改</el-button>
            <el-button type="danger" icon="Delete" text @click="()=>handleDelete([row.id])">删除</el-button>
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
import {SysRole} from "../../apis/sys_role.ts";
import {SysMenu} from "../../apis/sys_menu.ts";
import FormDialog from "../../components/common/FormDialog.vue";
const treeRef = ref()
const pageLoading = ref(true)
const menus = ref([])
const tableData = ref([])
const dialogOpen = ref(false)
const submitForm = ref({})

const getMenus = async ()=>{
  const response = await SysMenu.List()
  menus.value = response.data.list
}

const getPageData = async () => {
  pageLoading.value = true
  const response = await SysRole.List()
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
      treeRef.value.setCheckedKeys(defaultForm.editMenu.map(item=>item.id))
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
    console.log(menus)
    const response = await SysRole.UpdatePermission(submitForm.value.id,menus)
    console.log(response)
    ElMessage.success("权限修改成功")
  }else if (!submitForm.value.id){
    const response = await SysRole.Create({
      ...submitForm.value,
    })
    ElMessage.success("创建成功")
    console.log(response)
  }else{
    const response = await SysRole.Edit(submitForm.value)
    console.log(response)
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
          const response = await SysRole.Delete(ids)
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
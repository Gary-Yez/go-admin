<template>
  <el-card class="container" shadow="never" v-loading="pageLoading">
    <div class="mb-[10px] ">
      <el-button type="primary" :loading="pageLoading" @click="getPageData">刷新</el-button>
      <el-button type="primary" icon="Plus" @click="()=>handleAdd({})">新增管理员</el-button>
    </div>
    <el-table :data="tableData">
      <el-table-column label="编号" prop="id" :width="100"></el-table-column>
      <el-table-column label="头像" prop="avatar" :width="100">
        <template #default="{ row }">
          <el-avatar :src="row.avatar || '/img/user.png'"></el-avatar>
        </template>
      </el-table-column>
      <el-table-column label="昵称" prop="nickname"></el-table-column>
      <el-table-column label="用户名" prop="username"></el-table-column>
      <el-table-column label="角色" prop="role">
        <template #default="{ row }">
          <el-tag color="success">
            {{ rolesMap[row.role_id] ? rolesMap[row.role_id].name : '' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="手机号" prop="phone"></el-table-column>
      <el-table-column label="邮箱" prop="email"></el-table-column>
      <el-table-column label="状态" prop="status">
        <template #default="{ row }">
          <el-switch v-model="row.status" :active-value="1" :inactive-value="0"  :before-change="()=>handleChangeSwitch(row,'status')"></el-switch>
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
    <FormDialog v-model="dialogOpen" v-model:form="submitForm" :title="submitForm.id ? '修改管理员' : '新增管理员'" :on-confirm="handleSubmit">
      <el-form-item label="用户名" prop="username" :rules="[{required:true,message:'用户名不能为空'}]">
        <el-input v-model="submitForm.username" placeholder="请输入用户名"></el-input>
      </el-form-item>
      <el-form-item label="角色" prop="role_id" :rules="[{required:true,message:'请选择一个角色'}]">
        <el-select v-model="submitForm.role_id" filterable>
          <el-option v-for="role in roles" :value="role.id" :label="role.name" placeholder="请选择角色"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="昵称" prop="nickname" :rules="[{required:true,message:'昵称不能为空'}]">
        <el-input v-model="submitForm.nickname" placeholder="请输入昵称"></el-input>
      </el-form-item>
      <el-form-item label="手机号" prop="phone" :rules="[{required:true,message:'手机号不能为空'}]">
        <el-input v-model="submitForm.phone" placeholder="请输入手机号"></el-input>
      </el-form-item>
      <el-form-item label="邮箱" prop="email" :rules="[{required:true,message:'邮箱不能为空'}]">
        <el-input v-model="submitForm.email" placeholder="请输入邮箱"></el-input>
      </el-form-item>
      <el-form-item label="密码" prop="password" :rules="[{required:!submitForm.id,message:'密码不能为空'}]">
        <el-input v-model="submitForm.password" :placeholder="submitForm.id ? `无需重置密码可不填` : `请输入密码`"></el-input>
      </el-form-item>
    </FormDialog>
  </el-card>
</template>

<script setup lang="ts">
  import { SysAdminApi } from "../../apis/sys_admin.ts";
  import { SysRoleApi } from "../../apis/sys_role.ts";
  import {computed, onMounted, ref} from "vue";
  import { ElMessage,ElMessageBox } from "element-plus";
  import FormDialog from "../../../components/core/FormDialog.vue";

  const dialogOpen = ref(false)
  const roles:any = ref([])
  const queryForm = ref({
    page:1,
    limit:10
  })
  const pageLoading = ref(true)
  const total = ref(0)
  const tableData = ref([])
  const submitForm:any = ref({})


  const rolesMap = computed(()=>{
    let mapData:any = {}
    roles.value.forEach((item:any) => {
      mapData[item.id] = item
    })
    return mapData
  })

  const getRoles = async ()=>{
    const response = await SysRoleApi.List()
    roles.value = response.data.list
  }

  const getPageData = async () => {
    pageLoading.value = true
    try {
      const response = await SysAdminApi.List(queryForm.value)
      tableData.value = response.data.list
      total.value = response.data.total
    }catch (e) {
      console.log(e)
    }
    pageLoading.value = false
  }

  const handleAdd = (defaultForm:any)=>{
    submitForm.value = defaultForm
    dialogOpen.value = true
  }

  const handleDelete = (ids:Array<any>) => {
    ElMessageBox.confirm("您确认要删除该管理员吗？","删除提示",{
      type:"error",
      beforeClose:async (action:any, instance:any, done:any)=>{
        if (action === "confirm") {
          instance.confirmButtonLoading = true
          try {
            await SysAdminApi.Delete(ids)
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

  const handleChangeSwitch = async (row:any,key:string) => {
    row.loading = true
    try {
      await SysAdminApi.Edit({
        ...row,
        [key]:row[key] ? 0 : 1,
      })
      ElMessage.success("修改成功")
      row.loading = false
      return true
    }catch (e) {
      row.loading = false
      return false
    }
  }

  const handleSubmit = async ()=>{
    if (!submitForm.value.id){
      await SysAdminApi.Create(submitForm.value)
      ElMessage.success("创建成功")
    }else{
      await SysAdminApi.Edit(submitForm.value)
      ElMessage.success("修改成功")
    }
    getPageData()
  }

  onMounted(()=>{
    getRoles()
    getPageData()
  })
</script>


<style scoped>

</style>

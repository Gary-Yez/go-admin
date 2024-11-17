<template>
  <el-card class="container" shadow="never" v-loading="pageLoading">
    <div class="mb-[10px] ">
      <el-button type="primary" icon="Refresh" :loading="pageLoading" @click="getPageData">刷新</el-button>
      <el-button type="primary" icon="Plus" @click="()=>handleAdd({})">新增根菜单</el-button>
    </div>
    <el-table :data="tableData" row-key="id" default-expand-all>
      <el-table-column label="ID" prop="id" :width="100" sortable></el-table-column>
      <el-table-column label="字体图标" prop="icon" :width="100">
        <template #default="{ row }">
          <el-icon v-if="row.icon" :size="20">
            <component :is="row.icon"></component>
          </el-icon>
        </template>
      </el-table-column>
      <el-table-column label="菜单名称" prop="name"></el-table-column>
      <el-table-column label="路由地址" prop="path"></el-table-column>
      <el-table-column label="组件路径" prop="component"></el-table-column>
      <el-table-column label="父菜单ID" prop="parent_id" :width="100">
        <template #default="{ row }">
          <el-tag v-if="!row.parent_id" type="primary">根菜单</el-tag>
          <el-tag v-else type="success">{{ row.parent_id }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="排序" prop="sort" sortable></el-table-column>
      <el-table-column label="是否隐藏" prop="hidden" :width="100">
        <template #default="{ row }">
          <el-switch v-model="row.hidden" :loading="row.loading" :before-change="()=>handleChangeSwitch(row,'hidden')"></el-switch>
        </template>
      </el-table-column>
      <el-table-column label="操作" :width="280">
        <template #default="{ row }">
          <el-button-group class="table-btn-group">
            <el-button type="primary" icon="Plus" text @click="()=>handleAdd({parent_id:row.id})">添加子菜单</el-button>
            <el-button type="primary" icon="Edit" text @click="()=>handleAdd(row)">修改</el-button>
            <el-button type="danger" icon="Delete" text @click="()=>handleDelete([row.id])">删除</el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>
    <FormDialog v-model="dialogOpen" v-model:form="submitForm" :title="submitForm.id ? '修改菜单' : '新增菜单'" :on-confirm="handleSubmit">
      <el-form-item label="父菜单" prop="parent_id">
        <el-cascader class="w-full" v-model="submitForm.parent_id" :options="parentMenu" :props="{ label:'name',value:'id',checkStrictly:true,emitPath:false}"></el-cascader>
      </el-form-item>
      <el-row :gutter="15">
        <el-col :span="10">
          <el-form-item label="字体图标" prop="icon">
            <el-select v-model="submitForm.icon" placeholder="字体图标" filterable clearable>
              <template v-if="submitForm.icon" #prefix>
                <el-icon :size="16" class="mr-[5px]">
                  <component :is="submitForm.icon"></component>
                </el-icon>
              </template>
              <el-option v-for="icon in icons" :value="icon" :key="icon">
                <div class="flex items-center">
                  <el-icon :size="16" class="mr-[5px]">
                    <component :is="icon"></component>
                  </el-icon>
                  <span>{{ icon }}</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="14">
          <el-form-item label="菜单名称" prop="name" :rules="[{required:true,message:'请输入菜单名称'}]">
            <el-input v-model="submitForm.name" placeholder="请输入菜单名称"></el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item label="路由地址" prop="path" :rules="[{required:true,message:'请输入路由地址'}]">
        <el-input v-model="submitForm.path" placeholder="请输入路由地址">
          <template #prepend>/dashboard/</template>
        </el-input>
      </el-form-item>
      <el-row :gutter="12">
        <el-col :span="10">
          <el-form-item label="菜单类型">
            <el-radio-group v-model="submitForm.menu_type">
              <el-radio-button :value="1" label="根菜单"></el-radio-button>
              <el-radio-button :value="2" label="页面菜单"></el-radio-button>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="14" v-if="submitForm.menu_type === 2">
          <el-form-item label="组件路径" prop="component" :rules="[{required:true,message:'请选择组件'}]">
            <el-select  v-model="submitForm.component" placeholder="请选择组件" filterable>
              <el-option v-for="item in SyncComponents" :value="item" :label="item"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item label="菜单排序" prop="sort" :rules="[{required:true,message:'请输入菜单排序'}]">
        <el-input-number v-model="submitForm.sort" :value-on-clear="0"></el-input-number>
      </el-form-item>
    </FormDialog>
  </el-card>
</template>

<script setup lang="ts">
import {computed, onMounted, ref} from "vue";
import { SysMenu } from "../../apis/sys_menu.ts";
import FormDialog from "../../components/common/FormDialog.vue";
import {SyncComponents} from "../../routes/syncMenu.ts";

const pageLoading = ref(true)
const dialogOpen = ref(false);
const tableData = ref([])

const submitForm = ref({})
const icons = window.icons

const getPageData = async ()=>{
  pageLoading.value = true
  const response = await SysMenu.List()
  console.log(response)
  tableData.value = response.data.list
  pageLoading.value = false
}

const parentMenu = computed(()=>{
  return [{
    name:'根菜单',
    id:0,
    children:tableData.value
  }]
})

onMounted(()=>{
  getPageData()
})

const handleAdd = (defaultForm:any)=>{
  submitForm.value = {
    ...defaultForm,
    sort:defaultForm.sort || 0,
    menu_type: defaultForm.component ? 2 : 1,
    parent_id:defaultForm.parent_id || 0
  }
  dialogOpen.value = true
}

const handleDelete = (ids:Array<any>) => {
  ElMessageBox.confirm("您确认要删除该菜单吗？","删除提示",{
    type:"error",
    beforeClose:async (action, instance, done)=>{
      if (action === "confirm") {
        instance.confirmButtonLoading = true
        try {
          const response = await SysMenu.Delete(ids)
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

const handleChangeSwitch = async (row:any,key:string) => {
  row.loading = true
  try {
    await SysMenu.Edit({
      ...row,
      [key]:!row[key],
    })
    ElMessage.success("修改成功")
    row.loading = false
    return true
  }catch (e) {
    row.loading = false
    return false
  }
}

const handleSubmit = async () => {
  submitForm.value.parent_id = submitForm.value.parent_id || null
  if (!submitForm.value.id){
    const response = await SysMenu.Create(submitForm.value)
    ElMessage.success("创建成功")
  }else{
    const response = await SysMenu.Edit(submitForm.value)
    ElMessage.success("修改成功")
  }
  getPageData()
}

</script>

<style scoped>

</style>
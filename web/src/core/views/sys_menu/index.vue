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
          <iconify-icon class="text-[24px]" :icon="row.icon"></iconify-icon>
<!--          <el-icon v-if="row.icon" :size="20">-->
<!--            <component :is="row.icon"></component>-->
<!--          </el-icon>-->
        </template>
      </el-table-column>
      <el-table-column label="菜单名称" prop="name"></el-table-column>
      <el-table-column label="路由地址" prop="path"></el-table-column>
      <el-table-column label="组件路径" prop="component">
        <template #default="{ row }">
          <el-tag v-if="row.component" type="success">{{ row.component }}</el-tag>
          <el-tag v-else type="primary">路由组件</el-tag>
        </template>
      </el-table-column>
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
            <el-button type="primary" icon="Plus" :disabled="!!row.component" text @click="()=>handleAdd({parent_id:row.id})">添加子菜单</el-button>
            <el-button type="primary" icon="Edit" text @click="()=>handleAdd(row)">修改</el-button>
            <el-button type="danger" icon="Delete" text @click="()=>handleDelete([row.id])">删除</el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>
    <FormDialog v-model="dialogOpen" v-model:form="submitForm" :title="submitForm.id ? '修改菜单' : '新增菜单'" :on-confirm="handleSubmit">
      <el-form-item label="父菜单" prop="parent_id">
        <el-cascader class="w-full" v-model="submitForm.parent_id" :options="parentMenu" :props="{ label:'name',value:'id',checkStrictly:true,emitPath:false,disabled:'component'}"></el-cascader>
      </el-form-item>
      <el-row :gutter="15">
        <el-col :span="10">
          <el-form-item label="字体图标" prop="icon" :rules="[{required:true,message:'请选择字体图标'}]">
            <IconSelect v-model="submitForm.icon"></IconSelect>
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
          <template #prepend>{{ routePrefix }}</template>
        </el-input>
      </el-form-item>
      <el-row :gutter="12">
        <el-col :span="11">
          <el-form-item label="组件类型">
            <el-radio-group v-model="submitForm.menu_type">
              <el-radio-button :value="1" label="路由组件"></el-radio-button>
              <el-radio-button :value="2" label="页面组件"></el-radio-button>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="13" v-if="submitForm.menu_type === 2">
          <el-form-item label="组件路径" prop="component" :rules="[{required:true,message:'请选择组件'}]">
            <el-select v-model="submitForm.component" placeholder="请选择组件" filterable>
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
import { SysMenuApi } from "../../apis/sys_menu.ts";
import FormDialog from "../../../components/core/FormDialog.vue";
import {SyncComponents} from "../../../routes/syncMenu.ts";
import { ElMessage,ElMessageBox } from "element-plus";
import IconSelect from "../../../components/core/IconSelect.vue";

const pageLoading = ref(true)
const dialogOpen = ref(false);
const tableData = ref([])

const submitForm:any = ref({})


const getPageData = async ()=>{
  pageLoading.value = true
  const response = await SysMenuApi.List()
  tableData.value = response.data.list
  pageLoading.value = false
}

const parentMenu = computed(()=>{
  return [{
    name:'根菜单',
    id:0,
    children:tableData.value,
    path:"dashboard",
  }]
})

onMounted(()=>{
  getPageData()
})

function findParentNodes(root:any, targetId:any) {
  const result:any = [];
  const traverse = (node:any, path:any) => {
    if (node.id === targetId) {
      // 找到目标节点，将路径中的父级存入结果
      result.push(...[...path, node]);
      return true;
    }
    if (node.children) {
      // 继续遍历子节点，并传递当前路径（父级 + 当前节点）
      for (const child of node.children) {
        if (traverse(child, [...path, node])) {
          return true; // 找到后提前终止遍历
        }
      }
    }
    return false;
  };
  traverse(root[0], []); // 初始路径为空
  return result;
}

const routePrefix = computed(()=>{
  console.log(parentMenu.value,submitForm.value.parent_id)
  let treeNode = findParentNodes(parentMenu.value,submitForm.value.parent_id)
  return treeNode.reduce((pre:string,cur:any)=>{
    return pre + cur.path + "/";
  },"/")
})

const handleAdd = (defaultForm:any)=>{
  submitForm.value = {
    ...defaultForm,
    sort:defaultForm.sort || 0,
    menu_type: defaultForm.component ? 2 : 1,
    parent_id:defaultForm.parent_id || 0,
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
          await SysMenuApi.Delete(ids)
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
    await SysMenuApi.Edit({
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
  let form = JSON.parse(JSON.stringify(submitForm.value))
  form.parent_id = form.parent_id || null
  delete form.parents
  if (!form.id){
    await SysMenuApi.Create(form)
    ElMessage.success("创建成功")
  }else{
    await SysMenuApi.Edit(form)
    ElMessage.success("修改成功")
  }
  getPageData()
}

</script>

<style scoped>

</style>

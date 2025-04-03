<template>
  <div class="container" v-loading="pageLoading">
    <div v-if="isDev">
      <el-card header="模块信息" shadow="never">
        <el-form ref="formRef" :model="submitForm">
          <el-row :gutter="15">
            <el-col :span="6">
              <el-form-item label="模块名称" prop="module_name" :rules="[{required:true,message:'请输入模块名称'}]">
                <el-input v-model="submitForm.module_name" :disabled="isEdit" placeholder="请输入模块名称" @input="(val:any)=>submitForm.model_name = snakeToCamel(val)"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="中文名称" prop="chinese_module_name" :rules="[{required:true,message:'请输入中文名称'}]">
                <el-input v-model="submitForm.chinese_module_name" :disabled="isEdit" placeholder="请输入中文名称"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="结构体名称" prop="model_name" :rules="[{required:true,message:'结构体名称'}]">
                <el-input v-model="submitForm.model_name" disabled placeholder="根据模块名称自动生成"></el-input>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
        <div>
          <el-checkbox v-model="submitForm.use_common" disabled>使用标准数据结构</el-checkbox>
          <el-checkbox v-model="submitForm.create_curd" @change="handleChangeCURD">生成CURD</el-checkbox>
          <el-checkbox v-model="submitForm.use_soft_delete" :disabled="!submitForm.create_curd">使用软删除</el-checkbox>
        </div>
      </el-card>
      <el-card v-if="submitForm.create_curd" class="mt-[15px]" header="结构体字段" shadow="never">
        <el-form ref="fieldFormRef" :model="submitForm">
          <el-table :data="submitForm.fields">
            <el-table-column label="字段名称">
              <template #default="{ row,$index }">
                <el-form-item :prop="`fields.${$index}.name`" :rules="[{required:true,message:'请输入字段名称'}]">
                  <el-input v-model="row.name" placeholder="请输入字段名称" @input="(val:any)=>handleChangeFieldName(row,val)"></el-input>
                </el-form-item>
              </template>
            </el-table-column>
            <el-table-column label="字段类型">
              <template #default="{ row,$index }">
                <el-form-item :prop="`fields.${$index}.type`" :rules="[{required:true,message:'请选择字段类型'}]">
                  <el-select v-model="row.type" placeholder="字段类型">
                    <el-option label="字符串" value="string"></el-option>
                    <el-option label="布尔值" value="bool"></el-option>
                    <el-option label="整数型" value="int"></el-option>
                    <el-option label="无符号整数" value="uint"></el-option>
                    <el-option label="浮点数" value="float64"></el-option>
                  </el-select>
                </el-form-item>
              </template>
            </el-table-column>
            <el-table-column label="结构体Key">
              <template #default="{ row,$index }">
                <el-form-item :prop="`fields.${$index}.key`">
                  <el-input v-model="row.key" disabled placeholder="根据字段名称自动生成"></el-input>
                </el-form-item>
              </template>
            </el-table-column>
            <el-table-column label="中文名称">
              <template #default="{ row,$index }">
                <el-form-item :prop="`fields.${$index}.chinese_name`">
                  <el-input v-model="row.chinese_name" placeholder="请输入中文名称"></el-input>
                </el-form-item>
              </template>
            </el-table-column>
            <el-table-column label="表格可见" :width="80">
              <template #default="{ row,$index }">
                <el-form-item :prop="`fields.${$index}.table_show`">
                  <el-switch v-model="row.table_show" :disabled="row.hidden"></el-switch>
                </el-form-item>
              </template>
            </el-table-column>
            <el-table-column label="可编辑" :width="80">
              <template #default="{ row,$index }">
                <el-form-item :prop="`fields.${$index}.editable`">
                  <el-switch v-model="row.editable" :disabled="row.hidden"></el-switch>
                </el-form-item>
              </template>
            </el-table-column>
            <el-table-column label="是否必填" :width="80">
              <template #default="{ row,$index }">
                <el-form-item :prop="`fields.${$index}.required`">
                  <el-switch v-model="row.required" :disabled="row.hidden"></el-switch>
                </el-form-item>
              </template>
            </el-table-column>
            <el-table-column label="隐藏" :width="80">
              <template #default="{ row,$index }">
                <el-form-item :prop="`fields.${$index}.hidden`">
                  <el-switch v-model="row.hidden" @change="(val:any)=>{ row.table_show = !val;row.editable = !val;row.required = !val }"></el-switch>
                </el-form-item>
              </template>
            </el-table-column>
            <el-table-column label="索引类型">
              <template #default="{ row,$index }">
                <el-form-item :prop="`fields.${$index}.index_type`">
                  <el-select v-model="row.index_type" placeholder="索引类型" clearable>
                    <el-option label="index" value="index"></el-option>
                    <el-option label="unique" value="unique"></el-option>
                  </el-select>
                </el-form-item>
              </template>
            </el-table-column>
            <el-table-column label="操作" :width="100">
              <template #default="{ row }">
                <el-button class="mb-[18px]" size="small" type="danger" icon="Delete" :disabled="submitForm.fields.length <= 1" @click="()=>handleDeleteField(row.id)">删除</el-button>
              </template>
            </el-table-column>
            <template #append>
              <el-button class="w-full" type="primary" @click="handleAddFiled">新增字段</el-button>
            </template>
          </el-table>
        </el-form>
      </el-card>
      <el-card  class="mt-[15px]" shadow="never">
        <div class="text-right">
          <el-button type="primary" @click="handlePreview">预览代码</el-button>
        </div>
      </el-card>
      <el-dialog title="代码预览" class="preview-dialog" v-model="previewShow" :before-close="(done:any)=>!submitLoading && done()">
        <el-alert type="error" show-icon :closable="false">若待生成的文件已存在，则会覆盖文件内的所有内容并无法找回！</el-alert>
        <el-tabs class="mt-[10px]" tab-position="left">
          <el-tab-pane v-for="item in previewList" :label="item.path">
            <div class="pre-box">
              <pre v-text="item.content"></pre>
            </div>
          </el-tab-pane>
        </el-tabs>
        <template #footer>
          <div style="flex: auto">
            <el-button size="large" @click="!submitLoading && (previewShow = false)">取消</el-button>
            <el-button size="large" type="primary" :loading="submitLoading" @click="handleConfirm">确认生成</el-button>
          </div>
        </template>
      </el-dialog>
    </div>
    <el-card class="container"  v-else>
      <el-result title="此功能只允许开发环境使用" icon="error"></el-result>
    </el-card>
  </div>
</template>

<script setup lang="ts">
  import {onMounted, ref} from "vue";
  import {SysAutocodeApi} from "../../../apis/core/sys_autocode.ts";
  import {useRoute} from "vue-router";
  import { ElMessage } from "element-plus";

  const isDev = import.meta.env.MODE === 'development';
  const route = useRoute()
  const isEdit = ref(false)
  const initForm = ()=>{
    return {
      id:+new Date(),
      table_show:true,
      editable:true,
      required:true,
    }
  }
  const formRef = ref()
  const fieldFormRef = ref()
  const pageLoading = ref(true)
  const submitForm = ref({
    use_common:true,
    create_curd:true,
    use_soft_delete:false,
    chinese_module_name:"",
    module_name:"",
    model_name:"",
    fields:[initForm()]
  })
  const previewShow = ref(false)
  const previewList:any = ref([])
  const submitLoading = ref(false)


  function snakeToCamel(str:string) {
    return str.replace(/_([a-z])/g, (_:string, letter:string) => letter.toUpperCase())
        .replace(/^[a-z]/, (match:string) => match.toUpperCase());
  }

  const handleChangeFieldName = (row:any,val:string) => {
    row.key = snakeToCamel(val)
    row.json_name = val
  }

  const handleChangeCURD = () => {
    submitForm.value.fields = [initForm()]
    submitForm.value.use_soft_delete = false
  }

  const handleAddFiled = () => {
    submitForm.value.fields.push(initForm())
  }

  const handleDeleteField = (id:number) => {
    submitForm.value.fields = submitForm.value.fields.filter((item:any)=>{
      return item.id !== id
    })
  }

  const handlePreview = async () => {
    pageLoading.value = true
    try {
      await formRef.value.validate()
      if (fieldFormRef.value){
        await fieldFormRef.value.validate()
      }
      const response = await SysAutocodeApi.Preview(submitForm.value)
      previewList.value = response.data
      previewShow.value = true
    }catch (e) {
      console.log(e)
    }
    pageLoading.value = false
  }

  const handleConfirm = async () => {
    submitLoading.value = true
    try {
      await SysAutocodeApi.Generate(submitForm.value)
      previewShow.value = false
      ElMessage.success("生成成功")
    }catch (e) {
      console.log(e)
    }
    submitLoading.value = false
  }


  onMounted(async () => {
    if (route.query.id){
      isEdit.value = true
      pageLoading.value = true
      const response = await SysAutocodeApi.GetHistory(route.query.id)
      submitForm.value = JSON.parse(response.data.form)
      pageLoading.value = false
    }else{
      pageLoading.value = false
    }
  })


</script>

<style lang="less" scoped>
  :deep(.preview-dialog){
    width: 90%;
    .el-dialog__body{
      .el-tabs{
        height: 100%;
        .el-tabs__item{
          font-size: 14px;
          &.is-left{
            justify-content: flex-start;
          }
        }
        .el-tab-pane{
          height: 50vh;
          .pre-box{
            overflow: auto;
            height: 100%;
            padding: 15px;
            background: #334155;
            color: #ffffff;
          }
        }
      }
    }
  }
</style>

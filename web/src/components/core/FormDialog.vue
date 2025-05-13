<template>
  <el-drawer
      v-model="show"
      :before-close="handleClose"
      :close-on-click-modal="false"
  >
    <template #header>
      <el-page-header @back="handleClose">
        <template #content>
          <span>{{ props.title }}</span>
        </template>
      </el-page-header>
    </template>
    <el-form :model="form" label-position="top" ref="formRef" size="large">
      <slot></slot>
    </el-form>
    <template #footer>
      <el-button size="large" @click="handleClose">取消</el-button>
      <el-button size="large" type="primary" :loading="confirmLoading" @click="handleConfirm">确认</el-button>
    </template>
  </el-drawer>
<!--  <el-dialog class="max-w-[500px]" v-model="show" :title="props.title" :before-close="handleClose" :close-on-click-modal="false">-->

<!--  </el-dialog>-->
</template>

<script setup lang="ts">
import { ref } from "vue";

const formRef = ref();
const show = defineModel({
  default(){
    return false
  }
})

const form = defineModel("form",{
  default(){
    return {

    }
  }
})

const props = defineProps({
  title:{
    type:String,
  },
  onConfirm:{
    type:Function,
  }
})

const confirmLoading = ref(false)

const handleClose = (done?:any) => {
  if (confirmLoading.value) {
    return
  }
  show.value = false
  form.value = {}
  formRef.value.resetFields()
  if (typeof done === 'function') {
    done()
  }
}

const handleConfirm = async () => {
  await formRef.value.validate()
  confirmLoading.value = true
  try {
    props.onConfirm && await props.onConfirm()
    confirmLoading.value = false
    handleClose()
  }catch (e) {
    confirmLoading.value = false
    console.log(e)
  }
}
</script>


<style scoped lang="less">

</style>

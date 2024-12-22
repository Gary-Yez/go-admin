<template>
  <el-dialog class="max-w-[500px]" v-model="show" :title="props.title" :before-close="handleClose">
    <el-form :model="form" label-position="top" ref="formRef" size="large">
      <slot></slot>
    </el-form>
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="confirmLoading" @click="handleConfirm">确认</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import {ref} from "vue";

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

const handleClose = (done:any) => {
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
    show.value = false
  }catch (e) {
    console.log(e)
  }
  confirmLoading.value = false
}
</script>


<style scoped>

</style>
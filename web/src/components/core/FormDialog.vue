<template>
  <el-drawer
      class="w-[100%!important]"
      :style="`max-width: ${props.maxWidth}px`"
      v-model="show"
      :before-close="handleClose"
      :close-on-click-modal="props.closeOnClickModal"
      :close-on-press-escape="props.closeOnPressEscape"
      :destroy-on-close="props.destroyOnClose"
  >
    <template #header>
      <el-page-header @back="handleClose">
        <template #content>
          <span>{{ props.title }}</span>
        </template>
      </el-page-header>
    </template>
    <el-form :model="form" label-position="top" ref="formRef" :size="props.size">
      <slot :formRef="formRef"></slot>
    </el-form>
    <template #footer>
      <el-button size="large" @click="handleClose">{{ props.cancelBtnText }}</el-button>
      <el-button size="large" type="primary" :loading="confirmLoading" @click="handleConfirm">{{ props.confirmBtnText }}</el-button>
    </template>
  </el-drawer>
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
    return {}
  }
})
type SizeType = 'large' | 'default' | 'small'

const props = withDefaults(defineProps<{
  title?: string
  onConfirm?: () => Promise<void>
  closeOnClickModal?: boolean
  closeOnPressEscape?: boolean
  destroyOnClose?: boolean
  maxWidth?: number
  size?: SizeType
  confirmBtnText?:string
  cancelBtnText?:string
}>(), {
  closeOnClickModal: false,
  closeOnPressEscape: false,
  destroyOnClose: false,
  maxWidth: 500,
  size: "large",
  confirmBtnText:"确认",
  cancelBtnText:"取消"
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

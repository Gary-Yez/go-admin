<template>
  <template v-for="menu in props.menus">
    <el-menu-item v-if="!menu.hidden && !menu.children" :index="BuildMenuPath(menu)">
      <el-icon v-if="menu.icon">
        <iconify-icon :icon="menu.icon"></iconify-icon>
      </el-icon>
      <span>{{ menu.name }}</span>
    </el-menu-item>
    <el-sub-menu v-else-if="!menu.hidden" :index="BuildMenuPath(menu)">
      <template #title>
        <el-icon>
          <iconify-icon v-if="menu.icon" :icon="menu.icon"></iconify-icon>
        </el-icon>
        <span>{{ menu.name }}</span>
      </template>
      <MenuItem :prefix="BuildMenuPath(menu)" :menus="menu.children"></MenuItem>
    </el-sub-menu>
  </template>
</template>

<script setup lang="ts">
const props:any = defineProps({
  menus:{
    type:Array,
    default:()=>{
      return []
    }
  },
  prefix:{
    type:String,
    default:()=>{
      return "/dashboard"
    }
  }
})
const BuildMenuPath = (menu:any) => {
  return props.prefix + "/" + menu.path
}
</script>

<style scoped>

</style>

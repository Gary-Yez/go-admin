<template>
  <template v-for="menu in props.menus">
    <el-menu-item v-if="!menu.hidden && !menu.children" :index="getMenuPath(menu.path)">
      <el-icon v-if="menu.icon">
        <component :is="menu.icon"></component>
      </el-icon>
      <span>{{ menu.name }}</span>
    </el-menu-item>
    <el-sub-menu v-else-if="!menu.hidden" :index="getMenuPath(menu.path)">
      <template #title>
        <el-icon v-if="menu.icon">
          <component :is="menu.icon"></component>
        </el-icon>
        <span>{{ menu.name }}</span>
      </template>
      <MenuItem :menus="menu.children"></MenuItem>
    </el-sub-menu>
  </template>
</template>

<script setup lang="ts">
const props = defineProps({
  menus:{
    type:Array,
    default:()=>{
      return []
    }
  }
})

const getMenuPath = (path: string) => {
  return `/dashboard/${path}`
}
</script>

<style scoped>

</style>
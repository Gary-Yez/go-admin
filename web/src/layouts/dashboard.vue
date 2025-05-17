<template>
  <el-container class="main-layout">
    <el-header class="main-header">
      <div class="flex h-full justify-between items-center">
        <div class="flex items-center">
          <div class="logo">
            <img class="logo-img" src="/logo.png" alt="">
            <div class="logo-text">GoAdmin</div>
          </div>
          <div>
            <el-breadcrumb separator="/">
              <el-breadcrumb-item v-for="item in $route.matched">{{ item.meta.name }}</el-breadcrumb-item>
            </el-breadcrumb>
            <div class="text-[12px] text-[var(--el-text-color-regular)] mt-[5px]">当前时间：{{ commonStore.currentTime }}</div>
          </div>
        </div>
        <div class="flex items-center">
          <div class="theme-btn">
            <el-icon v-if="commonStore.theme === 'dark'" :size="18" @click="()=>commonStore.setTheme('light')">
              <iconify-icon icon="iconoir:sun-light"></iconify-icon>
            </el-icon>
            <el-icon v-else :size="18" @click="()=>commonStore.setTheme('dark')">
              <iconify-icon icon="iconoir:half-moon"></iconify-icon>
            </el-icon>
          </div>
          <el-dropdown trigger="hover" @command="handleCommand">
            <div class="flex items-center outline-none">
              <el-avatar class="mr-[5px]" :src="userStore.UserData?.avatar" :size="30"></el-avatar>
              <span class="text-[14px]">{{ userStore.UserData?.nickname }}</span>
              <el-icon class="el-icon--right">
                <iconify-icon icon="iconoir:nav-arrow-down-solid"></iconify-icon>
              </el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="">当前角色：{{ userStore.UserData?.role?.name  }}</el-dropdown-item>
                <el-dropdown-item command="userinfo">个人信息</el-dropdown-item>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>

        </div>
      </div>
    </el-header>
    <el-container class="content-layout">
      <el-aside class="content-aside">
        <div class="aside">
          <el-menu class="main-slider-menu" router :default-active="$route.path">
            <MenuItem :menus="menus"></MenuItem>
          </el-menu>
        </div>
      </el-aside>
      <el-main class="content-content">
        <el-tabs v-model="activeTab" class="nav-tabs" type="card"  @tab-remove="removeTab" @tab-click="onTabClick">
          <el-tab-pane v-for="tab in visitedTabs" :key="tab.fullPath" :name="tab.fullPath" :closable="visitedTabs.length > 1">
            <template #label>
              <el-dropdown trigger="contextmenu">
                <div class="tab-label">
                  <iconify-icon class="mr-[5px]" :icon="tab.meta.icon"></iconify-icon>
                  <span>{{ tab.meta.name }}</span>
                </div>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="()=>handleRemoveOtherTabs(tab.fullPath)">关闭其他</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </el-tab-pane>
        </el-tabs>
        <div class="dashboard-page">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
  import { onMounted, ref,watch} from "vue";
  import {useUserStore} from "../stores/user.ts";
  import MenuItem from "../components/core/MenuItem.vue";
  import {useCommonStore} from "../stores/common.ts";
  import {useRoute, useRouter} from "vue-router";
  const commonStore = useCommonStore()
  const userStore = useUserStore()
  const route = useRoute()
  const router = useRouter()
  const menus:any = ref([])
  commonStore.setTheme()
  commonStore.setTime()
  const activeTab = ref(route.path)
  const visitedTabs = ref([
    {
      fullPath: route.fullPath,
      meta: route.meta,
    },
  ])

  setInterval(()=>{
    commonStore.setTime()
  },1000)

  onMounted(()=>{
    menus.value = userStore.UserMenu
  })

  // 添加新标签页
  const addTab = (route:any) => {
    if (!visitedTabs.value.find((t) => t.fullPath === route.fullPath)) {
      visitedTabs.value.push({
        fullPath: route.fullPath,
        meta: route.meta,
      })
    }
    activeTab.value = route.fullPath
  }

  // 切换标签页时导航
  const onTabClick = (tab:any) => {
    router.push(tab.paneName)
  }

  // 关闭标签页
  const removeTab = (targetName:string) => {
    const tabs = visitedTabs.value
    const index = tabs.findIndex((t) => t.fullPath === targetName)
    if (tabs[index].fullPath === activeTab.value) {
      const nextTab = tabs[index + 1] || tabs[index - 1]
      if (nextTab) router.push(nextTab.fullPath)
    }
    visitedTabs.value = tabs.filter((t) => t.fullPath !== targetName)
  }
  //关闭其他
  const handleRemoveOtherTabs = (targetName:string) => {
    const tabs = visitedTabs.value
    visitedTabs.value = tabs.filter((t) => t.fullPath === targetName)
    router.push(targetName)
  }

  const handleCommand = (command:string) => {
    switch (command) {
      case "userinfo":
        router.push("/dashboard/sys_userinfo")
        break
      case "logout":
        userStore.logout()
        break
    }
  }

  // 监听路由变化添加标签
  watch(
      () => route.fullPath,
      (_newPath:string) => {
        addTab(route)
      },
      { immediate: true }
  )


</script>

<style lang="less" scoped>

.main-layout{
  width: 100%;
  height: 100%;
  background: var(--main-bg-light-color);
  .main-header{
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    z-index: 1;
    background: var(--main-bg-color);
    padding: 0 15px;
    //margin-bottom: 15px;
    box-shadow:0 1px 1px 1px rgba(0, 0, 0, .1);
    .logo{
      display: flex;
      align-items: center;
      margin-right: 30px;
      .logo-img{
        width: 30px;
        height: 30px;
        margin-right: 10px;
      }
      .logo-text{
        font-size: 21px;
        font-weight: bold;
        color: var(--main-text-color);
      }
    }
    .theme-btn{
      margin-right: 15px;
      --tw-shadow: 0 1px 3px 0 rgb(0 0 0 / .1), 0 1px 2px -1px rgb(0 0 0 / .1) !important;
      --tw-shadow-colored: 0 1px 3px 0 var(--tw-shadow-color), 0 1px 2px -1px var(--tw-shadow-color) !important;
      box-shadow: var(--tw-ring-offset-shadow, 0 0 #0000), var(--tw-ring-shadow, 0 0 #0000), var(--tw-shadow) !important;
      border-radius: 50%;
      border:1px solid var(--el-border-color-lighter);
      padding: 5px;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
    }
  }
  .content-layout{
    height: 100%;
    padding-top: 60px;
    //padding: 15px;
    .content-aside{
      overflow: hidden;
      //flex: unset !important;
      width: 256px !important;
      max-width: 256px !important;
      height: 100%;
      padding: 8px 8px;
      background: var(--main-bg-color);
      .aside{
        overflow: hidden;
        height: 100%;
        :deep(.main-slider-menu){
          width: 100%;
          height: 100%;
          border-right: none;
          .el-sub-menu{
            &.is-active{
              .el-sub-menu__title{
                color: var(--el-color-primary);
              }
            }
          }
          .el-menu-item{
            &.is-active{
              background: var(--el-color-primary);
              color: #ffffff;
            }
          }
        }
      }
    }
    .content-content{
      display: flex;
      flex-direction: column;
      padding: 0;
      .nav-tabs{
        background-color: var(--main-bg-color);
        :deep(.el-tabs__header){
          margin-bottom: 0;
          border: none;
          .el-tabs__nav{
            border-bottom: unset;
            border: 1px solid var(--el-border-color-light);
            .el-tabs__item{
              &.is-active{
                background: var(--main-bg-light-color);
                .tab-label{
                  color: var(--el-color-primary);
                }
              }
              .el-dropdown{
                height: 100%;
                .el-tooltip__trigger{
                  display: flex;
                  align-items: center;
                }
              }
            }
          }
        }
      }
      .dashboard-page{
        flex: 1;
        overflow: hidden;
        overflow-y: auto;
        padding: 12px;
      }
    }
  }
}
</style>

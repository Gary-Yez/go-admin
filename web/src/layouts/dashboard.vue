<template>
  <el-container class="main-layout">
    <el-aside class="main-aside">
      <div class="aside">
        <div class="logo">
          <img class="logo-img" src="/logo.svg" alt="">
          <div class="logo-text">管理员后台</div>
        </div>
        <el-menu class="main-slider-menu" router :default-active="$route.path">
          <MenuItem :menus="menus"></MenuItem>
          <MenuItem v-if="commonStore.isDev" :menus="DevMenu"></MenuItem>
        </el-menu>
      </div>
    </el-aside>
    <el-container class="content-layout">
      <el-header class="content-header">
        <div class="flex h-full justify-between items-center">
          <div>
            <el-breadcrumb separator="/">
              <el-breadcrumb-item v-for="item in $route.matched">{{ item.meta.name }}</el-breadcrumb-item>
            </el-breadcrumb>
            <div class="text-[12px] text-[var(--el-text-color-regular)] mt-[5px]">当前时间：{{ commonStore.currentTime }}</div>
          </div>
          <div class="flex items-center">
            <div class="theme-btn">
              <el-icon v-if="commonStore.theme === 'dark'" :size="18" @click="()=>commonStore.setTheme('light')">
                <component is="Sunny"></component>
              </el-icon>
              <el-icon v-else :size="18" @click="()=>commonStore.setTheme('dark')">
                <component is="Moon"></component>
              </el-icon>
            </div>
            <el-dropdown trigger="hover" @command="handleCommand">
              <div class="flex items-center outline-none">
                <el-avatar class="mr-[5px]" :src="userStore.UserData?.avatar" :size="30"></el-avatar>
                <span class="text-[14px]">{{ userStore.UserData?.nickname }}</span>
                <el-icon class="el-icon--right">
                  <component is="ArrowDown"></component>
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
      <el-main class="content-content">
        <div class="dashboard-page">
          <router-view v-slot="{ Component }">
            <transition name="fade-transform" mode="out-in" enter-from-class="fade-transform-enter">
              <component :is="Component"></component>
            </transition>
          </router-view>
        </div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
  import { onMounted, ref} from "vue";
  import {useUserStore} from "../stores/user.ts";
  import MenuItem from "../components/common/MenuItem.vue";
  import {useCommonStore} from "../stores/common.ts";
  import {DevMenu} from "../routes/syncMenu.ts";
  const commonStore = useCommonStore()
  const userStore = useUserStore()
  const menus:any = ref([])

  commonStore.setTheme()
  commonStore.setTime()

  setInterval(()=>{
    commonStore.setTime()
  },1000)

  onMounted(()=>{
    menus.value = userStore.UserMenu
  })

  const handleCommand = (command:string) => {
    switch (command) {
      case "userinfo":
        break
      case "logout":
        userStore.logout()
        break
    }
  }
</script>

<style lang="less" scoped>

.main-layout{
  width: 100%;
  height: 100%;
  background: var(--main-bg-light-color);
  .main-aside{
    overflow: hidden;
    flex: unset !important;
    width: 276px !important;
    max-width: 276px !important;
    height: 100%;
    padding: 15px 0 15px 15px;
    background: transparent;
    border-radius: 30px;
    .aside{
      overflow: hidden;
      border-radius: 15px;
      height: 100%;
      background: var(--main-bg-color);
      .logo{
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 150px;
        .logo-img{
          width: 50px;
          height: 50px;
          margin-bottom: 15px;
        }
        .logo-text{
          font-size: 24px;
          font-weight: bold;
          color: var(--main-text-color);
        }
      }
      :deep(.main-slider-menu){
        width: 100%;
        height: calc(100% - 150px);
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
  .content-layout{
    height: 100%;
    //padding: 15px;
    .content-header{
      background: var(--main-bg-color);
      padding: 0 15px;
      border-radius: 15px;
      margin: 15px;
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
    .content-content{
      overflow: hidden;
      overflow-y: auto;
      //padding: 15px 0 0 0;
      //margin: 0 15px 15px;
      padding: 0 15px 0;
      margin-bottom: 15px;
      .dashboard-page{
        height: 100%;
        //padding: 15px 15px 15px 0;
      }
    }
  }
}
</style>

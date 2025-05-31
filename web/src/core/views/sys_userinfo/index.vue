<template>
  <el-card shadow="never">
    <el-page-header v-loading="submitLoading" @back="$router.go(-1)">
      <template #content>个人信息</template>
      <div class="p-[15px]">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-card>
              <div class="text-center">
                <el-avatar :src="userStore.UserData.avatar || '/img/user.png'" :size="100"></el-avatar>
                <div class="text-[18px] my-[15px] font-bold">{{ userStore.UserData.username }}</div>
              </div>
              <el-descriptions title="详细信息" :column="1" border size="large">
                <el-descriptions-item label="昵称">{{ userStore.UserData.nickname }}</el-descriptions-item>
                <el-descriptions-item label="角色"><el-tag>{{ userStore.UserData.role?.name }}</el-tag></el-descriptions-item>
                <el-descriptions-item label="手机号">{{ userStore.UserData.phone }}</el-descriptions-item>
                <el-descriptions-item label="邮箱">{{ userStore.UserData.email }}</el-descriptions-item>
              </el-descriptions>
            </el-card>
          </el-col>
          <el-col :span="16">
            <el-card>
              <el-tabs @tab-change="handleChangeTab">
                <el-tab-pane label="基本信息" >
                  <el-form :model="userInfoForm" ref="userinfoRef" label-position="top" size="large">
                    <el-form-item label="昵称" prop="nickname" :rules="[{required:true,message:'请输入昵称'}]">
                      <el-input v-model="userInfoForm.nickname" placeholder="请输入昵称"></el-input>
                    </el-form-item>
                    <el-form-item label="手机号" prop="phone" :rules="[{required:true,message:'请输入手机号'}]">
                      <el-input v-model="userInfoForm.phone" placeholder="请输入手机号"></el-input>
                    </el-form-item>
                    <el-form-item label="邮箱" prop="email" :rules="[{required:true,message:'请输入邮箱'}]">
                      <el-input v-model="userInfoForm.email" placeholder="请输入邮箱"></el-input>
                    </el-form-item>
                    <div class="mt-[30px]">
                      <el-button type="primary" @click="handleChangeUserInfo">保存</el-button>
                    </div>
                  </el-form>
                </el-tab-pane>
                <el-tab-pane label="修改密码">
                  <el-form :model="passwordForm" ref="passwordRef" label-position="top" size="large">
                    <el-form-item label="旧密码" prop="old_password" :rules="[{required:true,message:'请输入旧密码'}]">
                      <el-input v-model="passwordForm.old_password" placeholder="请输入旧密码"></el-input>
                    </el-form-item>
                    <el-form-item label="新密码" prop="new_password" :rules="[{required:true,message:'请输入新密码'},{min:6,max:32,message: '密码长度必须要6-32位'}]">
                      <el-input v-model="passwordForm.new_password" placeholder="请输入新密码"></el-input>
                    </el-form-item>
                    <el-form-item label="确认密码" prop="confirm_password" :rules="[{required:true,message:'请确认密码'},{min:6,max:32,message: '密码长度必须要6-32位'}]">
                      <el-input v-model="passwordForm.confirm_password" placeholder="请确认密码"></el-input>
                    </el-form-item>
                    <div class="mt-[30px]">
                      <el-button type="primary" @click="handleChangePassword">确认修改</el-button>
                    </div>
                  </el-form>
                </el-tab-pane>
                <el-tab-pane label="API密钥">
                  <div v-if="userStore.UserData.api_token">
                    <el-input class="my-[20px]" v-model="userStore.UserData.api_token" size="large" disabled></el-input>
                    <el-row :gutter="20" class="mb-[15px]">
                      <el-col :span="12">
                        <el-button class="w-full" size="large" type="primary" :loading="submitLoading" @click="handleResetApiToken">重置API密钥</el-button>
                      </el-col>
                      <el-col :span="12">
                        <el-button class="w-full" size="large" type="success" :loading="submitLoading" @click="handleCopyApiToken">复制密钥</el-button>
                      </el-col>
                    </el-row>
                    <el-alert type="warning" :closable="false">重置密钥后，旧的API密钥有效时间为一小时</el-alert>
                  </div>
                  <div v-else>
                    <el-empty description="暂无API密钥，点击下方按钮生成">
                      <template #default>
                        <el-button size="large" type="primary" :loading="submitLoading" @click="handleResetApiToken">生成API密钥</el-button>
                      </template>
                    </el-empty>
                  </div>
                </el-tab-pane>
              </el-tabs>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-page-header>
  </el-card>
</template>

<script setup lang="ts">
  import {useUserStore} from "../../../stores/user.ts";
  import {ref} from "vue";
  import {SysAuthApi} from "../../apis/sys_auth.ts";
  import {ElMessage} from "element-plus";
  import {copyText} from "../../../utils/utils.ts";
  const submitLoading = ref(false);
  const userinfoRef = ref()
  const passwordRef = ref()
  const userStore = useUserStore()
  const initUserForm = ()=>{
    return {
      nickname:userStore.UserData.nickname,
      phone:userStore.UserData.phone,
      email:userStore.UserData.email,
    }
  }
  const initPasswordForm = ()=>{
    return {
      old_password:"",
      new_password:"",
      confirm_password:""
    }
  }
  const userInfoForm = ref(initUserForm())
  const passwordForm = ref(initPasswordForm())

  const handleChangeTab = ()=>{
    userInfoForm.value = initUserForm()
    passwordForm.value = initPasswordForm()
    userinfoRef.value.resetFields()
    passwordRef.value.resetFields()
  }

  const handleChangeUserInfo = async () => {
    await userinfoRef.value.validate()
    submitLoading.value = true;
    try {
      let res = await SysAuthApi.ChangeInfo(userInfoForm.value)
      ElMessage.success("修改成功")
      await userStore.getUserData()
      console.log(res)
    }catch (e) {
      console.log(e)
    }
    submitLoading.value = false;
  }

  const handleChangePassword = async () => {
    await passwordRef.value.validate()
    submitLoading.value = true;
    try {
      let res = await SysAuthApi.ChangePassword(passwordForm.value)
      console.log(res)
      ElMessage.success("修改成功")
      await userStore.logout()
    }catch (e) {
      console.log(e)
    }
    submitLoading.value = false;
  }

  const handleResetApiToken = async () => {
    submitLoading.value = true;
    try {
      const response = await SysAuthApi.ResetApiToken()
      userStore.UserData.api_token = response.data.api_token
      ElMessage.success("生成成功")
    }catch (e) {
      console.log(e)
    }
    submitLoading.value = false;
  }

  const handleCopyApiToken = () => {
    copyText(userStore.UserData.api_token || "")
  }
</script>

<style scoped lang="less">
.cell-item {
  display: flex;
  align-items: center;
}
</style>

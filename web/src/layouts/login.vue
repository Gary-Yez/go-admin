<template>
  <div class="container h-full">
   <div class="login-card">
     <div class="img-box">
       <img src="/img/login/bg.png" alt="">
     </div>
     <div class="form-box">
       <div class="logo-box">
         <img src="/logo.svg" alt="">
         <div class="logo-text">管理员登录</div>
       </div>
       <el-form :model="submitForm" ref="formRef" label-position="top" size="large">
         <el-form-item prop="username" :rules="[{ required:true,message:'请输入用户名' }]">
           <el-input v-model="submitForm.username" placeholder="请输入用户名">
             <template #prefix>
               <el-icon>
                 <component is="User"></component>
               </el-icon>
             </template>
           </el-input>
         </el-form-item>
         <el-form-item prop="username" :rules="[{ required:true,message:'请输入密码' }]">
           <el-input type="password" v-model="submitForm.password" placeholder="请输入密码">
             <template #prefix>
               <el-icon>
                 <component is="Lock"></component>
               </el-icon>
             </template>
           </el-input>
         </el-form-item>
         <div>
           <div class="remember-box">
             <el-checkbox v-model="isRemember">记住密码</el-checkbox>
           </div>
           <el-button class="logo-btn" type="primary" @click="handleSubmit">登录</el-button>
         </div>
       </el-form>
     </div>
   </div>
  </div>
</template>

<script setup lang="ts">
  import {ref} from "vue";
  import { SysAuthApi } from "../apis/sys_auth.ts";
  import {useUserStore} from "../stores/user.ts";
  import {useRouter} from "vue-router";

  const userStore = useUserStore()
  const router = useRouter()
  const formRef = ref()
  let rememberData = localStorage.getItem('remember')
  const submitForm = ref(rememberData ? JSON.parse(rememberData) :{
    username: "",
    password:""
  })
  const isRemember = ref(!!rememberData)


  const handleSubmit = async () => {
    await formRef.value.validate()
    const response = await SysAuthApi.Login(submitForm.value)
    userStore.setAccessToken(response.data.token)
    if (isRemember.value){
      localStorage.setItem("remember",JSON.stringify(submitForm.value))
    }else{
      localStorage.removeItem("remember")
    }
    await router.push('/dashboard')
  }
</script>


<style scoped>
  .container {
    background: linear-gradient(-135deg, #c850c0, #4158d0);
    .login-card{
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      display: flex;
      justify-content: center;
      background: #ffffff;
      border-radius: 20px;
      .img-box{
        width: 480px;
        display: flex;
        align-items: center;
        background: #0259e6;
        img{
          display: block;
          width: 100%;
        }
      }
      .form-box{
        width: 480px;
        padding: 100px 20px;
        .logo-box{
          text-align: center;
          margin-bottom: 40px;
          img{
            width: 50px;
            margin: 0 auto 10px;
          }
          .logo-text{
            font-size: 28px;
            font-weight: bold;
            color: var(--main-text-color);
          }
        }
        .remember-box{
          margin: 20px 0;
        }
        .logo-btn{
          width: 100%;
        }
      }
    }
  }
</style>

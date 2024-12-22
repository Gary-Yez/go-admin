import {defineStore} from "pinia";
import router from "../routes";
import {ElMessage} from "element-plus";

export const useUserStore = defineStore("userStore", {
    state: () => ({
        AccessToken: localStorage.getItem("access_token") || "",
        IsLogin: false,
        UserData:{} as {
            avatar?:string,
            nickname?:string,
            role?:{
                name:string,
                menus:[]
            }
        },
    }),
    actions: {
        setAccessToken(payload:string){
            this.AccessToken = payload
            localStorage.setItem("access_token", this.AccessToken)
        },
        setUserData(payload:Object){
            this.UserData = payload
            this.IsLogin = true
            ElMessage.success("已成功登录")
        },
        async logout(){
            localStorage.removeItem("access_token")
            this.AccessToken = ""
            this.IsLogin = false
            this.UserData = {}
            await router.replace("/login")
            ElMessage.success("已成功退出登录")
        }
    },
    getters:{
        UserMenu(state):Array<any>{
            return state.UserData?.role?.menus || []
        }
    }
})

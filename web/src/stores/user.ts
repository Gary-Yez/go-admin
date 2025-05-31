import {defineStore} from "pinia";
import router from "../routes";
import {ElMessage} from "element-plus";
import {SysAuthApi} from "../core/apis/sys_auth.ts";

export const useUserStore = defineStore("userStore", {
    state: () => ({
        AccessToken: localStorage.getItem("access_token") || "",
        IsLogin: false,
        UserData:{} as {
            avatar?:string,
            username?:string,
            nickname?:string,
            phone?:string,
            email?:string,
            api_token?:string,
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
        async getUserData(){
            const response = await SysAuthApi.GetMe()
            this.setUserData(response.data)
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
            let menus:any = state.UserData?.role?.menus || []
            if (import.meta.env.MODE === 'development'){
                menus.push({
                    name: "开发工具",
                    icon: "iconoir:laptop-dev-mode",
                    path: "sys_devtools",
                    children: [{
                        name:      "代码生成",
                        icon:      "iconoir:code",
                        path:      "autocode",
                        component: "../core/views/sys_devtools/autocode.vue",
                    }, {
                        name:      "生成历史",
                        icon:      "iconoir:code-brackets-square",
                        path:      "autocode_history",
                        component: "../core/views/sys_devtools/autocode_history.vue",
                    }],
                })
            }
            return menus
        }
    }
})

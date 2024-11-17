import {createRouter,createWebHistory} from 'vue-router';
import {useUserStore} from "../stores/user.ts";
import { SysAuth } from "../apis/sys_auth.ts";
import {addSyncRouter, getBaseRouter} from "./syncMenu.ts";
import NProgress from "nprogress"
const router = createRouter({
    history: createWebHistory(),
    routes:[{
        path:"/",
        redirect:"/login"
    },{
        path:"/login",
        name:'login',
        component:()=>import("../views/Login.vue"),
    },getBaseRouter()],
})


router.beforeEach(async (to, from, next) => {
    NProgress.start();
    const userStore = useUserStore()
    if (userStore.AccessToken && !userStore.IsLogin){
        try {
            const response = await SysAuth.GetMe()
            userStore.setUserData(response.data)
            addSyncRouter(userStore.UserMenu)
            return next(to.path)
        }catch (e) {
            await userStore.logout()
            return next("/login")
        }
    }
    if (!userStore.IsLogin && to.name !== "login"){
        return next('/login')
    }else if (userStore.IsLogin && to.name === "login"){
        return next('/dashboard')
    }
    next()
})

router.afterEach((to, from) => {
    NProgress.done();
})

export default router;
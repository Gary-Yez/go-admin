import {createRouter, createWebHashHistory} from 'vue-router';
import {useUserStore} from "../stores/user.ts";
import {addSyncRouter, getBaseRouter, layoutsModules} from "./syncMenu.ts";
//@ts-ignore
import NProgress from "nprogress"

const router = createRouter({
    history: createWebHashHistory(),
    routes:[{
        path:"/",
        redirect:"/login"
    },{
        path:"/login",
        name:'login',
        component:layoutsModules["../layouts/login.vue"],
    },getBaseRouter()],
})


router.beforeEach(async (to, _, next) => {
    NProgress.start();
    const userStore = useUserStore()
    if (userStore.AccessToken && !userStore.IsLogin){
        try {
            await userStore.getUserData()
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

router.afterEach(() => {
    NProgress.done();
})

export default router;

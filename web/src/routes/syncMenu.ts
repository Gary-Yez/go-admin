import router from "./index.ts";

const layoutsModules = import.meta.glob('../layouts/**/*.vue')
const viewModules = import.meta.glob('../views/**/*.vue')

const flatMenuTreeToRouter = (menus:Array<any>,routers:Array<any>)=>{
    if (!routers){
        routers = []
    }
    menus.forEach((item)=>{
        routers.push({
            path:`/dashboard/${item.path}`,
            meta:{
                name:item.name,
            },
            component:viewModules[item.component] || undefined,
        })
        if (item.children){
            routers = flatMenuTreeToRouter(item.children,routers)
        }
    })
    return routers
}

export const getBaseRouter = ()=>{
    return {
        path: '/dashboard',
        name:'dashboard',
        component:layoutsModules["../layouts/dashboard.vue"],
        meta:{
            name:"根目录",
        },
        children:[]
    }
}

export const addSyncRouter = (menus:Array<any>)=>{
    let routes:any = flatMenuTreeToRouter(menus,[])
    routes = routes.filter((item:any)=>item.component)
    const dashboardRouter = getBaseRouter()
    dashboardRouter.children = routes
    router.addRoute(dashboardRouter)
}

export const SyncComponents = Object.keys(viewModules)
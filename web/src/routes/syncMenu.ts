import router from "./index.ts";

const coreModules = import.meta.glob('../core/views/**/*.vue')
export const layoutsModules = import.meta.glob('../layouts/**/*.vue')
export const viewModules = import.meta.glob('../views/**/*.vue')

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

const flatMenuTreeToRouter = (menus:Array<any>,routers:Array<any>)=>{
    if (!routers){
        routers = []
    }
    menus.forEach((item)=>{
        let route:any = {
            path:item.path,
            meta:{
                name:item.name,
                icon:item.icon,
            },
        }
        if (viewModules[item.component] || coreModules[item.component]){
            route.component = viewModules[item.component] || coreModules[item.component]
        }
        if (item.children){
            route.children = flatMenuTreeToRouter(item.children,[])
        }
        routers.push(route)
    })
    return routers
}

export const addSyncRouter = (menus:Array<any>)=>{
    let routes:any = flatMenuTreeToRouter(menus,[])
    const dashboardRouter:any = getBaseRouter()
    dashboardRouter.children = routes
    if (routes.length > 0){
        dashboardRouter.redirect = `/dashboard/${routes[0].path}`
    }
    router.addRoute(dashboardRouter)
}

export const SyncComponents = Object.keys(viewModules)

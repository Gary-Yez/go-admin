import router from "./index.ts";

export const layoutsModules = import.meta.glob('../layouts/**/*.vue')
const autocodeModules = import.meta.glob('../autocode/**/*.vue')
export const viewModules = import.meta.glob('../views/**/*.vue')

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
            component:viewModules[item.component] || autocodeModules[item.component] || undefined,
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
export const DevMenu = [
    {
        name: "开发工具",
        icon: "EditPen",
        path: "dev",
        children: [{
            name:      "代码生成",
            icon:      "Cpu",
            path:      "autocode",
            component: "../autocode/index.vue",
        }, {
            name:      "生成历史",
            icon:      "Document",
            path:      "sys_autocode_history",
            component: "../autocode/history.vue",
        }],
    }
]

export const addSyncRouter = (menus:Array<any>)=>{
    let routes:any = flatMenuTreeToRouter(menus,[])
    routes = routes.filter((item:any)=>item.component)
    const dashboardRouter:any = getBaseRouter()
    if (import.meta.env.MODE === 'development'){
        routes = flatMenuTreeToRouter(DevMenu,routes)
        console.log("ok")
    }
    dashboardRouter.children = routes
    if (routes.length > 0){
        dashboardRouter.redirect = routes[0].path
    }

    router.addRoute(dashboardRouter)
}

export const SyncComponents = Object.keys(viewModules)

import { createApp } from 'vue'
import App from './App.vue'
import {createPinia} from "pinia";
import router from "./routes";
import ElementPlus from "element-plus";
import "./assets/css/tailwind.less"
import "nprogress/nprogress.css"
import "element-plus/dist/index.css"
import 'element-plus/theme-chalk/dark/css-vars.css'
import "./assets/css/global.less"
import "./assets/css/theme.less"
import * as ElementPlusIconsVue from "@element-plus/icons-vue"

const pinia = createPinia()
const app = createApp(App)

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
app.use(ElementPlus).use(pinia).use(router).mount('#app')

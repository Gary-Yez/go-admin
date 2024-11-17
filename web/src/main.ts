import { createApp } from 'vue'
import App from './App.vue'
import {createPinia} from "pinia";
import router from "./routes";
import "./assets/css/tailwind.less"
import "nprogress/nprogress.css"
import 'element-plus/theme-chalk/dark/css-vars.css'
import "./assets/css/global.less"
import "./assets/css/theme.less"
const pinia = createPinia()
const app = createApp(App)

import * as ElementPlusIconsVue from '@element-plus/icons-vue'
window.icons = []
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
    window.icons.push(key)
}
app.use(pinia).use(router).mount('#app')

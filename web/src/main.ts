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

const pinia = createPinia()
const app = createApp(App)


app.use(ElementPlus).use(pinia).use(router).mount('#app')

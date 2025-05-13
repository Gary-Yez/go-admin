import axios from "axios"
import {useUserStore} from "../stores/user.ts";
import {ElMessage} from "element-plus";

const request = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL,
})

request.interceptors.request.use(config=>{
    const userStore = useUserStore()
    config.headers.Authorization = `Bearer ${userStore.AccessToken}`
    return config
})


request.interceptors.response.use(response=>{
    if (response.data.code !== 200){
        ElMessage.error(response.data.message)
        return Promise.reject(response.data.message)
    }
    return response.data
},(error)=>{
    ElMessage.error(error.message)
    return Promise.reject(error)
})

export {
    request
}

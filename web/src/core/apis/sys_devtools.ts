import {request} from "../../utils/request.ts";

export const SysDevtoolsApi = {
    Preview(formData:any){
        return request.post("/sys_devtools/preview", formData)
    },
    Generate(formData:any){
        return request.post("/sys_devtools/generate", formData)
    },
    History(query:any){
        return request.get("/sys_devtools/history", {
            params:query
        })
    },
    GetHistory(id:any){
        return request.get("/sys_devtools/get_history", {
            params: {
                id
            }
        })
    },
    Delete(ids:Array<any>){
        return request.post("/sys_devtools/delete_history", {
            ids
        })
    }
}

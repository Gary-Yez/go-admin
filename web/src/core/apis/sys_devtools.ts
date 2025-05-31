import {request} from "../../utils/request.ts";

export const SysAutocodeApi = {
    Preview(formData:any){
        return request.post("/sys_autocode/preview", formData)
    },
    Generate(formData:any){
        return request.post("/sys_autocode/generate", formData)
    },
    History(query:any){
        return request.get("/sys_autocode/history", {
            params:query
        })
    },
    GetHistory(id:any){
        return request.get("/sys_autocode/get_history", {
            params: {
                id
            }
        })
    },
    Delete(ids:Array<any>){
        return request.post("/sys_autocode/delete_history", {
            ids
        })
    }
}

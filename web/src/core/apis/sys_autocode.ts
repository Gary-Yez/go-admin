import {request} from "../../utils/request.ts";

export const SysAutocodeApi = {
    Preview(formData:any){
        return request.post("/sys/autocode/preview", formData)
    },
    Generate(formData:any){
        return request.post("/sys/autocode/generate", formData)
    },
    History(query:any){
        return request.get("/sys/autocode/history", {
            params:query
        })
    },
    GetHistory(id:any){
        return request.get("/sys/autocode/get_history", {
            params: {
                id
            }
        })
    },
    Delete(ids:Array<any>){
        return request.post("/sys/autocode/delete_history", {
            ids
        })
    }
}

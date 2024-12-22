import {request} from "../utils/request.ts";

export const SysAutocodeApi = {
    Preview(formData:any){
        return request.post("/autocode/preview", formData)
    },
    Generate(formData:any){
        return request.post("/autocode/generate", formData)
    },
    History(query:any){
        return request.get("/autocode/history", {
            params:query
        })
    },
    GetHistory(id:any){
        return request.get("/autocode/get_history", {
            params: {
                id
            }
        })
    },
    Delete(ids:Array<any>){
        return request.post("/autocode/delete_history", {
            ids
        })
    }
}

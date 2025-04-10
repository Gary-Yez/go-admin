import { request } from "../../utils/request.ts";

export const SysTaskApi = {
    GetRegisteredHandler(){
        return request.get("/sys/task/get_registered_handler");
    },
    Get(id:number){
        return request.get("/sys/task/get",{
            params:{
                id:id
            }
        });
    },
    List(query:any){
        return request.post("/sys/task/list",{
            params:query
        });
    },
    Create(form:any){
        return request.post("/sys/task/create",{
            ...form,
            params:form.params ? JSON.stringify(form.params) : ""
        });
    },
    Edit(form:any){
        console.log(form.params)
        return request.post("/sys/task/edit", {
            ...form,
            params:form.params ? JSON.stringify(form.params) : ""
        });
    },
    Delete(ids:Array<number>){
        return request.post("/sys/task/delete", {
            ids:ids
        });
    }
}

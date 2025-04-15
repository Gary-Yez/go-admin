import { request } from "../../utils/request.ts";

export const SysCronApi = {
    GetRegisteredHandler(){
        return request.get("/sys/cron/get_registered_handler");
    },
    Get(id:number){
        return request.get("/sys/cron/get",{
            params:{
                id:id
            }
        });
    },
    List(query:any){
        return request.post("/sys/cron/list",{
            params:query
        });
    },
    Create(form:any){
        return request.post("/sys/cron/create",{
            ...form,
            params:form.params ? JSON.stringify(form.params) : ""
        });
    },
    Edit(form:any){
        console.log(form.params)
        return request.post("/sys/cron/edit", {
            ...form,
            params:form.params ? JSON.stringify(form.params) : ""
        });
    },
    Delete(ids:Array<number>){
        return request.post("/sys/cron/delete", {
            ids:ids
        });
    }
}

import {request} from "../../utils/request.ts";

export const SysGlobalVariableApi = {
    List(query:Object){
        return request.get("/sys/global_variable/list",{
            params:query
        });
    },
    Get(id:number){
        return request.get("/sys/global_variable/get",{
            params:{
                id:id
            }
        });
    },
    Create(form:Object){
        return request.post("/sys/global_variable/create",form);
    },
    Delete(ids:Array<number>){
        return request.post("/sys/global_variable/delete", {
            ids:ids
        });
    },
    Edit(form:Object){
        return request.post("/sys/global_variable/edit", form);
    }
}

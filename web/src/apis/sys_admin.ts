import {request} from "../utils/request.ts";

export const SysAdmin = {
    List(query){
        return request.get("/sys_admin/list",{
            params:query
        });
    },
    Create(form:any){
        return request.post("/sys_admin/create",form);
    },
    Delete(ids:Array<any>){
        return request.post("/sys_admin/delete", {
            ids:ids
        });
    },
    Edit(form:any){
        return request.post("/sys_admin/edit", form);
    }
}
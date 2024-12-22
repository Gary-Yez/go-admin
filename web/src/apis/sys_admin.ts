import {request} from "../utils/request.ts";

export const SysAdminApi = {
    List(query:Object){
        return request.get("/sys_admin/list",{
            params:query
        });
    },
    Create(form:Object){
        return request.post("/sys_admin/create",form);
    },
    Delete(ids:Array<number>){
        return request.post("/sys_admin/delete", {
            ids:ids
        });
    },
    Edit(form:Object){
        return request.post("/sys_admin/edit", form);
    }
}

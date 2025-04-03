import {request} from "../../utils/request.ts";

export const SysAdminApi = {
    List(query:Object){
        return request.get("/sys/admin/list",{
            params:query
        });
    },
    Create(form:Object){
        return request.post("/sys/admin/create",form);
    },
    Delete(ids:Array<number>){
        return request.post("/sys/admin/delete", {
            ids:ids
        });
    },
    Edit(form:Object){
        return request.post("/sys/admin/edit", form);
    }
}

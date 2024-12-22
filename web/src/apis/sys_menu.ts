import { request } from "../utils/request.ts";

export const SysMenuApi = {
    List() {
        return request.get("/sys_menu/list");
    },
    Create(form:Object){
        return request.post("/sys_menu/create",form);
    },
    Edit(form:Object){
        return request.post("/sys_menu/edit",form);
    },
    Delete(ids:Array<number>){
        return request.post("/sys_menu/delete", {
            ids:ids
        });
    }
}

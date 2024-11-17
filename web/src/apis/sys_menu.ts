import { request } from "../utils/request.ts";

export const SysMenu = {
    List() {
        return request.get("/sys_menu/list");
    },
    Create(form){
        return request.post("/sys_menu/create",form);
    },
    Edit(form){
        return request.post("/sys_menu/edit",form);
    },
    Delete(ids:Array<any>){
        return request.post("/sys_menu/delete", {
            ids:ids
        });
    }
}
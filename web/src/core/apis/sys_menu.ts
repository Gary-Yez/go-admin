import { request } from "../../utils/request.ts";

export const SysMenuApi = {
    List() {
        return request.get("/sys/menu/list");
    },
    Create(form:Object){
        return request.post("/sys/menu/create",form);
    },
    Edit(form:Object){
        return request.post("/sys/menu/edit",form);
    },
    Delete(ids:Array<number>){
        return request.post("/sys/menu/delete", {
            ids:ids
        });
    }
}

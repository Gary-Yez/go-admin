import { request } from "../../utils/request.ts";

export const SysApisApi = {
    SyncApi(){
        return request.get("/sys_apis/sync_api");
    },
    GetGroups(){
        return request.get("/sys_apis/get_groups");
    },
    UpdateIgnore(path:string,method:string,ignore:boolean){
        return request.post("/sys_apis/update_ignore", {
            path: path,
            method: method,
            ignore: ignore
        });
    },
    Get(id:number){
        return request.get("/sys_apis/get",{
            params:{
                id:id
            }
        });
    },
    List(query?:any){
        return request.post("/sys_apis/list",query);
    },
    Create(form:any){
        return request.post("/sys_apis/create",form);
    },
    Edit(form:any){
        return request.post("/sys_apis/edit", form);
    },
    Delete(ids:Array<number>){
        return request.post("/sys_apis/delete", {
            ids:ids
        });
    }
}

import {request} from "../utils/request.ts";

export const SysRoleApi  = {
    List(){
        return request.get("/sys_role/list");
    },
    Create(formData:any){
        return request.post("/sys_role/create",formData);
    },
    Edit(formData:any){
        return request.post("/sys_role/edit",formData);
    },
    Delete(ids:Array<any>){
        return request.post("/sys_role/delete",{
            ids
        });
    },
    UpdatePermission(roleId:number,menuIds:Array<any>){
        return request.post("/sys_role/permission",{
            id:roleId,
            menus:menuIds.map(item=>({
                id:item
            }))
        })
    }
}

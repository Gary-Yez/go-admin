import {request} from "../../utils/request.ts";

export const SysRoleApi  = {
    List(){
        return request.get("/sys/role/list");
    },
    Create(formData:any){
        return request.post("/sys/role/create",formData);
    },
    Edit(formData:any){
        return request.post("/sys/role/edit",formData);
    },
    Delete(ids:Array<any>){
        return request.post("/sys/role/delete",{
            ids
        });
    },
    UpdatePermission(roleId:number,menuIds:Array<any>){
        return request.post("/sys/role/permission",{
            id:roleId,
            menus:menuIds.map(item=>({
                id:item
            }))
        })
    }
}

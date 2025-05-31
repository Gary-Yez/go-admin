import {request} from "../../utils/request.ts";

export const SysRoleApi  = {
    Get(id:number){
        return request.get("/sys_role/get",{
            params:{
                id:id
            }
        });
    },
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
    UpdatePermission(roleId:number,menuIds:Array<any>,apis:Array<any>){
        return request.post("/sys_role/permission",{
            id:roleId,
            menus:menuIds.map(item=>({
                id:item
            })),
            apis:apis,
        })
    }
}

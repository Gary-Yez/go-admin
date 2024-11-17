import {request} from "../utils/request.ts";

export const Test = {
    List(query){
        return request.get("/test/list",{
            params:query
        });
    },
    Get(id:number){
        return request.get("/test/get",{
            params:{
                id:id
            }
        });
    },
    Create(form:any){
        return request.post("/test/create",form);
    },
    Delete(ids:Array<number>){
        return request.post("/test/delete", {
            ids:ids
        });
    },
    Edit(form:any){
        return request.post("/test/edit", form);
    }
}
import { request } from "../utils/request.ts";

export const TestApi = {
    Get(id:number){
        return request.get("/test/get",{
            params:{
                id:id
            }
        });
    },
    List(query:any){
        return request.post("/test/list",query);
    },
    Create(form:any){
        return request.post("/test/create",form);
    },
    Edit(form:any){
        return request.post("/test/edit", form);
    },
    Delete(ids:Array<number>){
        return request.post("/test/delete", {
            ids:ids
        });
    }
}
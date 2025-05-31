import {request} from "../../utils/request.ts";

export const SysAuthApi = {
    Login(form:Object){
        return request.post("/sys_auth/login",form);
    },
    GetMe(){
        return request.get("/sys_auth/me");
    },
    ChangeInfo(form:any){
        return request.post("/sys_auth/change_info",form);
    },
    ChangePassword(form:any){
        return request.post("/sys_auth/change_password",form);
    }
}

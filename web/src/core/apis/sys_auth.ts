import {request} from "../../utils/request.ts";

export const SysAuthApi = {
    Login(form:Object){
        return request.post("/sys/auth/login",form);
    },
    GetMe(){
        return request.get("/sys/auth/me");
    },
    ChangeInfo(form:any){
        return request.post("/sys/auth/change_info",form);
    },
    ChangePassword(form:any){
        return request.post("/sys/auth/change_password",form);
    }
}

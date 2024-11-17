import {request} from "../utils/request.ts";

export const SysAuth = {
    Login(form){
        return request.post("/sys_auth/login",form);
    },
    GetMe(){
        return request.get("/sys_auth/me");
    }
}
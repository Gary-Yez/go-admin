import {request} from "../utils/request.ts";

export const SysAuthApi = {
    Login(form:Object){
        return request.post("/sys_auth/login",form);
    },
    GetMe(){
        return request.get("/sys_auth/me");
    }
}

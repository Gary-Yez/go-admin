import {request} from "../../utils/request.ts";

export const SysAuthApi = {
    Login(form:Object){
        return request.post("/sys/auth/login",form);
    },
    GetMe(){
        return request.get("/sys/auth/me");
    }
}

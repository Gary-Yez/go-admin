import { request } from "../../utils/request.ts";

export const SysHomeApi = {
    Statistic(){
        return request.get("/sys/home/statistic");
    }
}

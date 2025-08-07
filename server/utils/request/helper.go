package request

import "github.com/gin-gonic/gin"

func GetReq(ctx *gin.Context) (*Req, error) {
	var req = new(Req)
	err := ctx.ShouldBind(req)
	return req, err
}
func GetReqIds(ctx *gin.Context) (*ReqIds, error) {
	var ids = new(ReqIds)
	err := ctx.ShouldBind(ids)
	return ids, err
}

func GetReqList(ctx *gin.Context) (*ReqList, error) {
	var list = new(ReqList)
	err := ctx.ShouldBind(list)
	return list, err
}

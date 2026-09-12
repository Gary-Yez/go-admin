package sys_config

import (
	"github.com/Gary-Yez/go-admin/internal/bizconfig"
	"github.com/Gary-Yez/go-admin/response"
	"github.com/gin-gonic/gin"
)

type controllerStruct struct{}

func (*controllerStruct) SyncCache(ctx *gin.Context) {
	count, err := bizconfig.SyncCache(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, gin.H{"count": count})
}

func (*controllerStruct) CleanupInvalid(ctx *gin.Context) {
	body := new(CleanupBody)
	if err := ctx.ShouldBindJSON(body); err != nil {
		response.Error(ctx, err.Error())
		return
	}
	keys := make([]string, 0, len(body.Items))
	for _, item := range body.Items {
		keys = append(keys, item.Key)
	}
	count, err := bizconfig.CleanupInvalid(ctx.Request.Context(), keys)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, gin.H{"count": count})
}

func (*controllerStruct) Values(ctx *gin.Context) {
	list, err := bizconfig.Values(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.List(ctx, list, int64(len(list)))
}

func (*controllerStruct) UpdateValue(ctx *gin.Context) { changeValue(ctx, false) }
func (*controllerStruct) ResetValue(ctx *gin.Context)  { changeValue(ctx, true) }

func changeValue(ctx *gin.Context, reset bool) {
	body := new(ValueBody)
	if err := ctx.ShouldBindJSON(body); err != nil {
		response.Error(ctx, err.Error())
		return
	}
	if err := bizconfig.UpdateValue(ctx.Request.Context(), body.Key, body.Value, reset); err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, "配置已更新")
}

func (*controllerStruct) List(ctx *gin.Context) {
	definition, err := Service.List()
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, definition)
}

func (*controllerStruct) Preview(ctx *gin.Context) {
	body := new(SaveBody)
	if err := ctx.ShouldBindJSON(body); err != nil {
		response.Error(ctx, err.Error())
		return
	}
	preview, err := Service.Preview(body)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, preview)
}

func (*controllerStruct) Apply(ctx *gin.Context) {
	body := new(SaveBody)
	if err := ctx.ShouldBindJSON(body); err != nil {
		response.Error(ctx, err.Error())
		return
	}
	if err := Service.Apply(body); err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, "配置定义已更新")
}

func (*controllerStruct) Site(ctx *gin.Context) {
	ctx.Header("Cache-Control", "no-store")
	info, err := Service.Site()
	if err != nil {
		response.Error(ctx, "站点信息暂时不可用")
		return
	}
	response.Success(ctx, info)
}

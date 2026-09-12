package sys_monitor

import (
	"context"
	"github.com/Gary-Yez/go-admin/response"
	"github.com/gin-gonic/gin"
	"time"
)

type controllerStruct struct{}

func (*controllerStruct) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	rows, err := Service.List(ctx)
	if err != nil {
		response.Error(c, "节点监控读取失败，请稍后重试")
		return
	}
	response.Success(c, gin.H{"list": rows})
}

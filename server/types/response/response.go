package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(ctx *gin.Context, data ...any) {
	if len(data) == 0 {
		ctx.JSON(http.StatusOK, result{
			Code:    200,
			Message: "success",
		})
	} else {
		if str, ok := data[0].(string); ok {
			ctx.JSON(http.StatusOK, result{
				Code:    200,
				Message: str,
			})
		} else {
			ctx.JSON(http.StatusOK, result{
				Code:    200,
				Message: "success",
				Data:    data[0],
			})
		}
	}
	ctx.Abort()
}

func List(ctx *gin.Context, list interface{}, total int64) {
	ctx.JSON(http.StatusOK, result{
		Code: 200,
		Data: gin.H{
			"list":  list,
			"total": total,
		},
		Message: "success",
	})
}

func Error(ctx *gin.Context, errInfo any, errorCode ...int) {
	statusCode := 500
	if len(errorCode) != 0 {
		statusCode = errorCode[0]
	}
	if err, ok := errInfo.(error); ok {
		ctx.JSON(http.StatusOK, result{
			Code:    statusCode,
			Message: err.Error(),
		})
	} else if errString, ok := errInfo.(string); ok {
		ctx.JSON(http.StatusOK, result{
			Code:    statusCode,
			Message: errString,
		})
	} else {
		ctx.JSON(http.StatusOK, result{
			Code:    statusCode,
			Message: "未知错误",
		})
	}
	ctx.Abort()
}

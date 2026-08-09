package admin

import "github.com/gin-gonic/gin"

// Module is the only contract an application module needs to implement.
type Module interface {
	Name() string
	Initialize() error
	AdminRouter(adminGroup *gin.RouterGroup)
	PublicRouter(publicGroup *gin.RouterGroup)
}

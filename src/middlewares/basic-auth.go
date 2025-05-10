package middlewares

import (
	"app/config"

	"github.com/gin-gonic/gin"
)

var basicAuth gin.HandlerFunc

func BasicAuthMiddleware(c *gin.Context) {
	if config.Config.Server.Username == "" || config.Config.Server.Password == "" {
		c.Next()
		return
	}
	if basicAuth == nil {
		basicAuth = gin.BasicAuth(gin.Accounts{
			config.Config.Server.Username: config.Config.Server.Password,
		})
	}
	basicAuth(c)
}
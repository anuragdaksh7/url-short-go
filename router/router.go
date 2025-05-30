package router

import (
	"github.com/anuragdaksh7/url-short-go/internal/user"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler) {
	r = gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	r.POST("/user/create", userHandler.CreateUser)
	r.POST("/user/login", userHandler.LoginUser)
}

func Start(addr string) error {
	return r.Run(addr)
}

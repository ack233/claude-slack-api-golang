package routes

import (
	"net/http"
	"slackapi/app/middlewares"

	"github.com/gin-gonic/gin"
)

func Load(r *gin.Engine) {

	r.LoadHTMLGlob("app/templates/*")
	// 无权限路由组
	noAuthRouter := r.Group("/")
	{
		noAuthRouter.Any("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"title": "Main website",
			})
		})

		noAuthRouter.GET(("/login"), login)
		noAuthRouter.GET(("/callback"), callback)

	}

	authRouter := r.Group("/").Use(middlewares.JWTAuth())
	{

		authRouter.GET(("/backend-api/revoke"), revoke)
		authRouter.POST(("/backend-api/conversation"), middlewares.Check_notMethod("POST"), conversation)
	}
}

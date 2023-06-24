package appcmd

import (
	"slackapi/app/middlewares"
	"slackapi/app/routes"
	"slackapi/pkgs/config"
	"slackapi/pkgs/zlog"

	"github.com/gin-gonic/gin"
)

func Start() {

	listenport := config.ViperConfig.GetString("port")
	initGin := initRouter()
	zlog.SugLog.Infof("******服务初始化完成,监听端口为: %v******", listenport)
	err := initGin.Run("0.0.0.0:" + listenport)
	zlog.Fatalerror(err)
}

func initRouter() *gin.Engine {

	//设置运行模式
	var ginmode string = config.ViperConfig.GetString("ginmode")
	if ginmode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if ginmode != "debug" {
		zlog.SugLog.Fatalf("运行级别ginmode设置错误，无法识别%v", ginmode)
	}
	// 初始化引擎
	r := gin.New()

	// 公共中间件
	r.Use(middlewares.GinLogger())

	r.Use(middlewares.GinRecovery())

	r.Use(middlewares.CORSMiddleware())
	routes.Load(r)
	return r
}

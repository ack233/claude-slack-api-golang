package main

import (
	"slackapi/app/appcmd"
	"slackapi/pkgs/config"
	"slackapi/pkgs/initfunc"
	"slackapi/pkgs/zlog"
	"slackapi/slack"
)

func main() {
	//解析配置
	config.Init()

	//初始化日志系统
	zlog.Init()

	//初始化功能函数
	initfunc.InitFun()

	//启动slack websocket
	go slack.Run()

	//启动gin
	appcmd.Start()
}

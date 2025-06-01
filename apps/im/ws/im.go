package main

import (
	"flag"
	"fmt"
	"github/lhh-gh/easy-chat/apps/im/ws/internal/config"
	"github/lhh-gh/easy-chat/apps/im/ws/internal/handler"
	"github/lhh-gh/easy-chat/apps/im/ws/internal/svc"
	"github/lhh-gh/easy-chat/apps/im/ws/websocket"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

// main 是WebSocket服务器应用程序的入口函数
// 主要执行以下任务：
// 1. 解析命令行参数获取配置文件路径
// 2. 加载并验证服务器配置
// 3. 初始化服务上下文和带JWT认证的WebSocket服务器
// 4. 注册请求处理器并启动WebSocket服务器
func main() {
	// 解析命令行参数，获取配置参数
	flag.Parse()

	// 从配置文件加载服务器配置
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 初始化配置，若失败则panic
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	// 使用加载的配置创建服务上下文
	ctx := svc.NewServiceContext(c)

	// 初始化带JWT认证的WebSocket服务器
	// websocket.WithServerMaxConnectionIdle(10*time.Second) 可配置连接空闲超时
	srv := websocket.NewServer(c.ListenOn, websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)))
	defer srv.Stop() // 确保main函数退出时停止服务器

	// 向服务器注册所有请求处理器
	handler.RegisterHandlers(srv, ctx)

	// 启动WebSocket服务器并打印监听地址
	fmt.Println("WebSocket服务器启动于", c.ListenOn, "......")
	srv.Start()
}

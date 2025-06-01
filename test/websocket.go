package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrader 是websocket的连接升级器，采用默认配置
var upgrader = websocket.Upgrader{}

// serverWs 处理websocket请求，实现简单的消息回传功能
// 参数:
//
//	w http.ResponseWriter - HTTP响应写入器
//	r *http.Request - HTTP请求对象
//
// 功能说明:
//  1. 将HTTP连接升级为WebSocket连接
//  2. 持续读取客户端消息并原样返回
//  3. 发生错误时关闭连接
func serverWs(w http.ResponseWriter, r *http.Request) {
	// 升级HTTP连接到WebSocket协议
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("连接升级失败:", err)
		return
	}
	defer conn.Close() // 确保最终关闭连接

	// 消息处理主循环
	for {
		// 读取客户端消息
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Print("消息读取失败:", err)
			return
		}
		log.Printf("收到消息：%s", message)

		// 将消息原样返回给客户端
		response := fmt.Sprintf("服务器返回：%s", message)
		if err = conn.WriteMessage(messageType, []byte(response)); err != nil {
			log.Print("消息发送失败:", err)
			return
		}
	}
}

// main 函数设置WebSocket服务器
// 功能说明:
//  1. 注册/ws路径的WebSocket处理器
//  2. 启动2345端口监听服务
func main() {
	http.HandleFunc("/ws", serverWs)
	fmt.Println("WebSocket服务启动中...")
	log.Fatal(http.ListenAndServe("0.0.0.0:2345", nil))
}

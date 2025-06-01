package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Conn 表示一个WebSocket连接
type Conn struct {
	idleMu sync.Mutex // 互斥锁，用于保护空闲时间字段的并发访问

	Uid string // 用户唯一标识

	*websocket.Conn         // 嵌入的gorilla/websocket连接
	s               *Server // 所属的Server实例

	idle              time.Time     // 记录上次活动时间
	maxConnectionIdle time.Duration // 最大空闲时间

	done chan struct{} // 用于通知连接关闭的信号通道
}

// NewConn 创建并初始化一个新的WebSocket连接
func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	// 升级HTTP连接到WebSocket
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade err %v", err)
		return nil
	}

	// 初始化连接结构体
	conn := &Conn{
		Conn:              c,
		s:                 s,
		idle:              time.Now(),              // 设置初始活动时间
		maxConnectionIdle: s.opt.maxConnectionIdle, // 从服务器配置获取最大空闲时间
		done:              make(chan struct{}),     // 初始化关闭通道
	}

	// 启动心跳检测协程
	go conn.keepalive()

	return conn
}

// ReadMessage 读取消息并更新活动时间
func (c *Conn) ReadMessage() (messageType int, data []byte, err error) {
	messageType, data, err = c.Conn.ReadMessage()

	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.idle = time.Time{} // 重置为空时间表示连接处于活跃状态

	return
}

// WriteMessage 发送消息并更新活动时间
func (c *Conn) WriteMessage(messageType int, data []byte) error {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	// 注意：此方法不是线程安全的
	err := c.Conn.WriteMessage(messageType, data)
	c.idle = time.Now() // 更新最后活动时间为当前时间
	return err
}

// Close 关闭连接
func (c *Conn) Close() error {
	select {
	case <-c.done: // 如果已经关闭
	default:
		close(c.done) // 关闭通知通道
	}

	return c.Conn.Close() // 关闭底层WebSocket连接
}

// keepalive 心跳检测机制
func (c *Conn) keepalive() {
	// 创建定时器，初始超时时间为最大空闲时间
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer func() {
		idleTimer.Stop() // 确保定时器被停止
	}()

	for {
		select {
		case <-idleTimer.C: // 定时器触发
			c.idleMu.Lock()
			idle := c.idle
			if idle.IsZero() { // 如果连接处于活跃状态
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle) // 重置定时器
				continue
			}
			// 计算剩余超时时间
			val := c.maxConnectionIdle - time.Since(idle)
			c.idleMu.Unlock()
			if val <= 0 {
				// 已超时，关闭连接
				c.s.Close(c)
				return
			}
			idleTimer.Reset(val) // 重置定时器为剩余时间
		case <-c.done: // 收到关闭信号
			fmt.Println("客户端结束连接")
			return
		}
	}
}

//该代码实现的心跳检测机制是 空闲超时检测机制，主要特点和实现方式如下：
//机制类型：
//被动式心跳检测（非主动Ping/Pong）
//基于连接活动状态的超时管理
//核心实现逻辑：
//通过记录最后活动时间（idle字段）：
//ReadMessage/WriteMessage 操作会更新活动时间
//读取时设置 idle = time.Time{}（零值表示活跃状态）
//写入时设置 idle = time.Now()
//定时器循环检测（keepalive方法）：

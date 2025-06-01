package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type Server struct {
	sync.RWMutex

	opt            *serverOption
	authentication Authentication
	pattern        string

	routes map[string]HandlerFunc
	addr   string

	connToUser map[*Conn]string
	userToConn map[string]*Conn

	upgrader websocket.Upgrader
	logx.Logger
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)

	return &Server{
		routes:   make(map[string]HandlerFunc),
		addr:     addr,
		upgrader: websocket.Upgrader{},

		opt:            &opt,
		authentication: opt.Authentication,
		pattern:        opt.pattern,

		connToUser: make(map[*Conn]string),
		userToConn: make(map[string]*Conn),

		Logger: logx.WithContext(context.Background()),
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("server handler ws recorver err %v", r)
		}
	}()

	conn := NewConn(s, w, r)
	if conn == nil {
		return
	}
	/**
	// http升级为websocket
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade err %v", err)
	}
	**/

	// 鉴权
	if !s.authentication.Auth(w, r) {
		s.Send(&Message{Data: fmt.Sprint("不具备访问权限")}, conn)
		conn.Close()
		return
	}

	// 添加连接记录，会有并发问题
	s.addConn(conn, r)

	// 根据连接对象获取请求信息并处理
	go s.handlerConn(conn)
}

func (s *Server) addConn(conn *Conn, req *http.Request) {
	uid := s.authentication.UserId(req)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 验证用户是否之前登入过
	if c := s.userToConn[uid]; c != nil {
		// 关闭之前的连接
		c.Close()
	}

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

func (s *Server) GetUsers(conns ...*Conn) []string {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

func (s *Server) Close(conn *Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	if uid == "" {
		// 已经被关闭
		return
	}

	delete(s.connToUser, conn)
	delete(s.userToConn, uid)

	conn.Close()
}

// 根据连接对象执行任务处理
func (s *Server) handlerConn(conn *Conn) {
	uids := s.GetUsers(conn)
	conn.Uid = uids[0]

	for {
		// 获取请求消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message err %v", err)
			// 关闭并删除连接
			s.Close(conn)
			return
		}

		// 请求信息
		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			// s.Send(NewErrMessage(err), conn)
			// continue
			s.Errorf("json unmarshal err %v,msg %v", err, string(msg))
			s.Close(conn)
			return
		}

		// 依据请求消息类型分类处理
		switch message.FrameType {
		case FramePing:
			// ping：回复
			s.Send(&Message{FrameType: FramePing}, conn)
		case FrameData:
			// 根据请求的method分发路由并执行
			if handler, ok := s.routes[message.Method]; ok {
				handler(s, conn, &message)
			} else {
				// conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("不存在执行的方法：%v 请检查", message.Method)))
				s.Send(&Message{
					FrameType: FrameData,
					Data:      fmt.Sprintf("不存在执行的方法：%v 请检查", message.Method),
				}, conn)
			}
		}
	}
}

func (s *Server) SendByUserId(msg any, sendToUIds ...string) error {
	if len(sendToUIds) == 0 {
		return nil
	}

	return s.Send(msg, s.GetConns(sendToUIds...)...)
}

func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}

	// 编码
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

func (s *Server) Start() {
	http.HandleFunc(s.opt.pattern, s.ServerWs)
	fmt.Println("启动websocket")
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	fmt.Println("停止服务")
}

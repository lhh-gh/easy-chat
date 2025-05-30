package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github/lhh-gh/easy-chat/apps/user/api/internal/config"
	"github/lhh-gh/easy-chat/apps/user/rpc/userclient"
)

// ServiceContext 定义了服务上下文结构，包含配置和用户服务客户端
type ServiceContext struct {
	Config config.Config // 服务配置信息

	userclient.User // 用户服务客户端接口
}

// NewServiceContext 创建一个新的服务上下文实例
//
// 参数:
//   - c: 服务配置，包含用户RPC客户端配置等信息
//
// 返回值:
//   - *ServiceContext: 初始化后的服务上下文指针，包含配置和用户服务客户端实例
func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化用户服务RPC客户端并创建服务上下文
	return &ServiceContext{
		Config: c,

		User: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}

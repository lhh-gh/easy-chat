package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	// 增加UserRpc
	UserRpc zrpc.RpcClientConf

	JwtAuth struct {
		AccessSecret string
		//AccessExpire int64
	}
}

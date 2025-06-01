package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config 定义了服务的配置结构体
// 包含以下配置项：
//   - RpcServerConf: 继承zrpc的RPC服务基础配置
//   - Mysql: MySQL数据库配置
//   - DataSource: 数据库连接字符串
//   - Cache: 缓存配置，使用cache包的CacheConf类型
//   - Jwt: JWT认证配置
//   - AccessSecret: JWT签名密钥
//   - AccessExpire: JWT过期时间(秒)
type Config struct {
	zrpc.RpcServerConf
	Redisx redis.RedisConf
	Mysql  struct {
		DataSource string
	}

	Cache cache.CacheConf

	Jwt struct {
		AccessSecret string
		AccessExpire int64
	}
}

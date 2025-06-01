package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github/lhh-gh/easy-chat/apps/user/models"
	"github/lhh-gh/easy-chat/apps/user/rpc/internal/config"
	"github/lhh-gh/easy-chat/pkg/constants"
	"github/lhh-gh/easy-chat/pkg/ctxdata"
	"time"
)

type ServiceContext struct {
	Config config.Config
	*redis.Redis
	UserModels models.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:     c,
		Redis:      redis.MustNewRedis(c.Redisx),
		UserModels: models.NewUsersModel(sqlConn, c.Cache),
	}
}

// 设置超级权限的token
func (svc *ServiceContext) SetRootToken() error {

	systemToken, err := ctxdata.GetJwtToken(svc.Config.Jwt.AccessSecret, time.Now().Unix(), 99999, constants.SYSTEM_ROOT_UID)

	if err != nil {
		return err
	}
	return svc.Redis.Set(constants.REDIS_SYSTEM_ROOT_TOKEN, systemToken)
}

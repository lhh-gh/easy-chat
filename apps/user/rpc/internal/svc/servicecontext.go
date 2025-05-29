package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github/lhh-gh/easy-chat/apps/user/models"
	"github/lhh-gh/easy-chat/apps/user/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config

	UserModels models.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,

		UserModels: models.NewUsersModel(sqlConn, c.Cache),
	}
}

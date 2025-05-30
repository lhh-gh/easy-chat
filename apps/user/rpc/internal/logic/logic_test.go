package logic

import (
	"github.com/zeromicro/go-zero/core/conf"
	"github/lhh-gh/easy-chat/apps/user/rpc/internal/config"
	"github/lhh-gh/easy-chat/apps/user/rpc/internal/svc"
	"path/filepath"
)

var svcCtx *svc.ServiceContext

func init() {
	var c config.Config
	conf.MustLoad(filepath.Join("../../etc/dev/user.yaml"), &c)
	svcCtx = svc.NewServiceContext(c)
}

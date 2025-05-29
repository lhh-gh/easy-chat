package logic

import (
	"context"
	"github.com/pkg/errors"
	"github/lhh-gh/easy-chat/apps/user/models"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"github/lhh-gh/easy-chat/apps/user/rpc/internal/svc"
	"github/lhh-gh/easy-chat/apps/user/rpc/user"
)

// ErrUserNotFound 定义用户未找到的错误
var ErrUserNotFound = errors.New("没有这个用户")

// GetUserInfoLogic 包含获取用户信息的业务逻辑
type GetUserInfoLogic struct {
	ctx         context.Context     // 上下文
	svcCtx      *svc.ServiceContext // 服务上下文
	logx.Logger                     // 日志记录器
}

// NewGetUserInfoLogic 创建并返回一个新的GetUserInfoLogic实例
// 参数:
//   - ctx: 上下文对象
//   - svcCtx: 服务上下文对象
//
// 返回值:
//   - *GetUserInfoLogic: 新创建的GetUserInfoLogic指针
func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserInfo 根据用户ID获取用户信息
// 参数:
//   - in: 包含用户ID的请求参数
//
// 返回值:
//   - *user.GetUserInfoResp: 包含用户信息的响应
//   - error: 错误信息，如果用户不存在返回ErrUserNotFound
func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	// 从数据库查询用户信息
	userEntity, err := l.svcCtx.UserModels.FindOne(l.ctx, in.Id)

	if err != nil {
		// 处理查询错误
		if err == models.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 将数据库实体转换为响应实体
	var resp user.UserEntity
	// 使用Copier库将数据库实体转换为响应实体
	copier.Copy(&resp, userEntity)

	return &user.GetUserInfoResp{
		User: &resp,
	}, nil
}

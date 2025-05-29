package logic

import (
	"context"
	"github.com/pkg/errors"
	"github/lhh-gh/easy-chat/apps/user/models"
	"github/lhh-gh/easy-chat/pkg/ctxdata"
	"github/lhh-gh/easy-chat/pkg/encrypt"
	"time"

	"github/lhh-gh/easy-chat/apps/user/rpc/internal/svc"
	"github/lhh-gh/easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegister = errors.New("手机号没有注册")
	ErrUserPwdError     = errors.New("密码不正确")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// 1.验证用户是否注册，根据手机号验证
	userEntity, err := l.svcCtx.UserModels.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		if err != models.ErrNotFound {
			return nil, ErrPhoneNotRegister
		}

		return nil, err
	}

	// 2.验证密码
	if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
		return nil, ErrUserPwdError
	}

	// 3.生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, err
	}

	// 4.返回
	return &user.LoginResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}

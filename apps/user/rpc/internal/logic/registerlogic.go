package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github/lhh-gh/easy-chat/apps/user/models"
	"github/lhh-gh/easy-chat/pkg/ctxdata"
	"github/lhh-gh/easy-chat/pkg/encrypt"
	"github/lhh-gh/easy-chat/pkg/wuid"
	"time"

	"github/lhh-gh/easy-chat/apps/user/rpc/internal/svc"
	"github/lhh-gh/easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneIsRegister = errors.New("手机号已经注册过")
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// 1.验证用户是否注册，根据手机号验证
	userEntity, err := l.svcCtx.UserModels.FindByPhone(l.ctx, in.Phone)
	if err != nil && err != models.ErrNotFound {
		return nil, err
	}

	if userEntity != nil {
		return nil, ErrPhoneIsRegister
	}

	// 2.定义用户数据
	userEntity = &models.Users{
		Id:       wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Sex: sql.NullInt64{
			Int64: int64(in.Sex),
			Valid: true,
		},
	}

	// 3.密码加密
	if len(in.Password) > 0 {
		genPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, err
		}
		userEntity.Password = sql.NullString{
			String: string(genPassword),
			Valid:  true,
		}
	}

	// 4.插入数据
	_, err = l.svcCtx.UserModels.Insert(l.ctx, userEntity)
	if err != nil {
		return nil, err
	}

	// 5.生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, err
	}

	// 5.返回注册结果
	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}

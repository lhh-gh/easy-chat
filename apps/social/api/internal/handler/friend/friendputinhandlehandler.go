package friend

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github/lhh-gh/easy-chat/apps/social/api/internal/logic/friend"
	"github/lhh-gh/easy-chat/apps/social/api/internal/svc"
	"github/lhh-gh/easy-chat/apps/social/api/internal/types"
)

// 好友申请处理
func FriendPutInHandleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FriendPutInHandleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := friend.NewFriendPutInHandleLogic(r.Context(), svcCtx)
		resp, err := l.FriendPutInHandle(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

package check

import (
	"github.com/apicat/apicat/backend/i18n"
	"github.com/apicat/apicat/backend/model/user"
	"github.com/apicat/apicat/backend/module/auth"
	"github.com/apicat/apicat/backend/route/proto"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func checkMemberStatus(authorization string) *user.Users {
	if authorization == "" {
		return nil
	}

	parts := strings.SplitN(authorization, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil
	}

	mc, err := auth.ParseToken(parts[1])
	if err != nil {
		return nil
	}

	if mc.UserID == 0 {
		return nil
	}

	u, err := user.NewUsers(mc.UserID)
	if err != nil {
		return nil
	}

	return u
}

func CheckMember() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authorization := ctx.Request.Header.Get("Authorization")
		u := checkMemberStatus(authorization)

		if u == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    proto.InvalidOrIncorrectLoginToken,
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Auth.TokenParsingFailed"}),
			})
			ctx.Abort()
			return
		}

		if u == nil || u.IsEnabled == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    proto.InvalidOrIncorrectLoginToken,
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Auth.AccountDisabled"}),
			})
			ctx.Abort()
			return
		}

		//将当前请求的username信息保存到请求的上下文c上
		ctx.Set("CurrentUser", u)
		//后续的处理函数可以通过c.Get("CurrentUser")来获取请求的用户信息
		ctx.Next()
	}
}

package middleware

import (
	"strings"

	"sunflower-gin/api"

	"sunflower-gin/pkg/jwt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	tokenPrefix = "Bearer "

	CtxKeyUserID = "userId" // 用户ID上下文 key

)

// Auth 基于 JWT token 认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 token
		authorizationValue := c.GetHeader("Authorization")
		if len(authorizationValue) == 0 || !strings.HasPrefix(authorizationValue, tokenPrefix) {
			api.ResponseError(c, api.CodeNeedLogin)
			c.Abort()
			return
		}
		if len(authorizationValue) <= 7 || !strings.HasPrefix(authorizationValue, "Bearer ") {
			api.ResponseError(c, api.CodeInvalidToken)
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authorizationValue, "Bearer ")
		claims, err := jwt.ParseAccessToken(tokenString)
		if err != nil {
			zap.L().Sugar().Debugf("parse access token error: %v", err)
			api.ResponseError(c, api.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将用户ID存入上下文，后续中间件或业务逻辑可以直接从上下文中获取
		c.Set(CtxKeyUserID, claims.UserId)
		c.Next()
	}
}

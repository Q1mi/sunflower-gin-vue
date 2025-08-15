package jwt

import "github.com/golang-jwt/jwt/v5"

// tokenType 令牌类型
type tokenType string

const (
	accessToken  tokenType = "accessToken"
	refreshToken tokenType = "refreshToken"
)

// CustomClaims 自定义声明结构体并内嵌 jwt.RegisteredClaims
// jwt包自带的 jwt.RegisteredClaims 包含了官方字段
// 我们这里需要额外记录一个 username 字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
// JWTClaims 自定义声明结构体并内嵌 jwt.RegisteredClaims
type CustomClaims struct {
	UserId   int64  `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

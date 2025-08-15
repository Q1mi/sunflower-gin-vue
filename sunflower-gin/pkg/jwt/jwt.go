package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var obj *JWT

var (
	ErrInvalidTokenType = errors.New("invalid token type")
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("expired token")
)

type JWT struct {
	accessSecret         []byte // 访问令牌密钥
	refreshSecret        []byte // 刷新令牌密钥
	accessExpireSeconds  int64  // 访问令牌过期时间
	refreshExpireSeconds int64  // 刷新令牌过期时间
}

func NewJWT(viper *viper.Viper) *JWT {
	return &JWT{
		accessSecret:         []byte(viper.GetString("jwt.access_secret")),
		refreshSecret:        []byte(viper.GetString("jwt.refresh_secret")),
		accessExpireSeconds:  viper.GetInt64("jwt.access_expire_seconds"),
		refreshExpireSeconds: viper.GetInt64("jwt.refresh_expire_seconds"),
	}
}

func MustInit(cfg *viper.Viper) {
	obj = NewJWT(cfg)
}

// GenAccessToken 生成 access token
func GenAccessToken(userId int64, username string) (string, error) {
	return obj.genToken(userId, username, accessToken)
}

// GenRefreshToken 生成 refresh token
func GenRefreshToken(userId int64, username string) (string, error) {
	return obj.genToken(userId, username, refreshToken)
}

// genToken 生成token
func (j *JWT) genToken(userId int64, username string, typ tokenType) (string, error) {
	var (
		expiresAt time.Time
		secret    []byte
	)
	switch typ {
	case accessToken:
		expiresAt = time.Now().Add(time.Duration(j.accessExpireSeconds) * time.Second)
		secret = j.accessSecret
	case refreshToken:
		expiresAt = time.Now().Add(time.Duration(j.refreshExpireSeconds) * time.Second)
		secret = j.refreshSecret
	default:
		return "", ErrInvalidTokenType
	}
	zap.L().Sugar().Debugf("-->生成 %s token，过期时间：%v", typ, expiresAt)
	claims := &CustomClaims{
		UserId:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "liwenzhou.com",
			Subject:   "sunflower",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt), // 有效期
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := accessToken.SignedString(secret) // 签名
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// parseToken 解析 token
func (j *JWT) parseToken(tokenString string, typ tokenType) (*CustomClaims, error) {
	var claim CustomClaims
	token, err := jwt.ParseWithClaims(tokenString, &claim,
		func(token *jwt.Token) (interface{}, error) {
			switch typ {
			case accessToken:
				return j.accessSecret, nil
			case refreshToken:
				return j.refreshSecret, nil
			default:
				return nil, ErrInvalidTokenType
			}
		})
	if err != nil {
		return nil, err
	}

	if token.Valid { // 校验token
		return &claim, nil
	}
	return nil, ErrInvalidToken
}

// ParseAccessToken 解析 access token
func ParseAccessToken(tokenString string) (*CustomClaims, error) {
	return obj.parseToken(tokenString, accessToken)
}

// ParseRefreshToken 解析 refresh token
func ParseRefreshToken(tokenString string) (*CustomClaims, error) {
	return obj.parseToken(tokenString, refreshToken)
}

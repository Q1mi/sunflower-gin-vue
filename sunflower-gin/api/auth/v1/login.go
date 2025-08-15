package v1

// 登录请求的参数和响应结构体

type LoginReq struct {
	Username string `json:"username" binding:"required"` // 用户名，必填项
	Password string `json:"password" binding:"required"` // 密码，必填项
}

type LoginResp struct {
	AccessToken  string `json:"accessToken"`  // 登录成功后返回的 token
	RefreshToken string `json:"refreshToken"` // 登录成功后返回的 refresh token
}

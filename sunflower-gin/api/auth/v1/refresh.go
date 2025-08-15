package v1

type RefreshReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"` // 刷新令牌
}

type RefreshResp struct {
	AccessToken  string `json:"accessToken"`  // 访问令牌
	RefreshToken string `json:"refreshToken"` // 刷新令牌
}

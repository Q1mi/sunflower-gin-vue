package v1

// CreateReq 创建请求结构体
type CreateReq struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"eqfield=Password"`
}

// CreateRes 创建响应结构体
type CreateRes struct {
	UserId   int64  `json:"userId"`
	Username string `json:"username"`
}

type MeRes struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

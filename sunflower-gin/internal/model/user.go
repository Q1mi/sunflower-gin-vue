package model

type CreateUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type CreateUserOutput struct {
	UserId   int64  `json:"userId"`
	Username string `json:"username"`
}

type UserProfileOutput struct {
	UserId   int64  `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

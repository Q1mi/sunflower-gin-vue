package model

type LoginInput struct {
	Username string
	Password string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
}

type RefreshTokenOutput struct {
	AccessToken  string
	RefreshToken string
}

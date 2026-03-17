package auth

type LoginInput struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

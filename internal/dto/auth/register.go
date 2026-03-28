package auth

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email,max=256"`
	Username string `json:"username" validate:"required,min=2,max=30"`
	Password string `json:"password" validate:"required,min=8,max=30"`
	Timezone string `json:"timezone"`
	Language string `json:"language"`
	Theme    string `json:"theme"`
	Role     string `json:"user_role"`
}

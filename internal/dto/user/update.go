package user

type UpdateUserProfile struct {
	Username *string `json:"username"`
	Timezone *string `json:"timezone"`
	Language *string `json:"language"`
	Theme    *string `json:"theme"`
}

type UpdateUserPassword struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=30"`
}

type UpdateUserEmail struct {
	Password string `json:"password" validate:"required,email,max=256"`
	NewEmail string `json:"new_email" validate:"required,min=8,max=30"`
}

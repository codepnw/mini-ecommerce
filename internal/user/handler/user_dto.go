package userhandler

type UserCreateReq struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required,min=2"`
	LastName  string `json:"last_name" binding:"required,min=2"`
}

type UserUpdateReq struct {
	FirstName *string `json:"first_name,omitempty" binding:"omitempty,min=2"`
	LastName  *string `json:"last_name,omitempty" binding:"omitempty,min=2"`
}

type UserLoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

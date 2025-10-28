package userhandler

type UserCreateReq struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
}

type UserUpdateReq struct {
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=2"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=2"`
}

type UserLoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

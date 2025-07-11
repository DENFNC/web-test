package request

type AuthUserRequest struct {
	Login    string `json:"login" validate:"required,min=8,alphanum"`
	Password string `json:"password" validate:"required,min=8,password"`
}

func (req *AuthUserRequest) Validate() error {
	return validate.Struct(req)
}

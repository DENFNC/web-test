package request

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate = func() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		return len(s) >= 8 &&
			regexp.MustCompile(`[A-Z]`).MatchString(s) &&
			regexp.MustCompile(`[a-z]`).MatchString(s) &&
			regexp.MustCompile(`[0-9]`).MatchString(s) &&
			regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(s)
	})
	return v
}()

type RegisterUserRequest struct {
	Token    string `json:"token,omitempty"`
	Login    string `json:"login" validate:"required,min=8,alphanum"`
	Password string `json:"password" validate:"required,min=8,password"`
}

func (req *RegisterUserRequest) Validate() error {
	return validate.Struct(req)
}

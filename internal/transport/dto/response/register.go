package response

type RegisterUserResponse struct {
	Login string `json:"login"`
}

func (resp *RegisterUserResponse) Validate() error { return nil }

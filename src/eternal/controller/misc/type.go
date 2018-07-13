package misc

type SendSignupCodeRequest struct {
	PhoneNumber string `json:"phone_number" form:"phone_number" query:"phone_number" validate:"required"`
}

type SendSignupCodeResult struct {
	Wait int `json:"wait"`
}

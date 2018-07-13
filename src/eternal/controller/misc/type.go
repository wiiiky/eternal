package misc

type SendSignupCodeRequest struct {
	PhoneNumber string `json:"phone_number" form:"phone_number" query:"phone_number" validate:"required"`
}

type SendSignupCodeResult struct {
	Sent bool `json:"sent"`
	Wait int  `json:"wait"`
}

package account

type LoginRequest struct {
	CountryCode string `json:"country_code" form:"country_code" query:"country_code"`
	PhoneNumber string `json:"phone_number" form:"phone_number" query:"phone_number" validate:"required"`
	Password    string `json:"password" form:"password" query:"password" validate:"required"`
}

type SignupRequest struct {
	LoginRequest
	Code string `json:"code" form:"code" query:"code" validate:"required"`
}

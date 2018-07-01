package account

type LoginRequest struct {
	CountryCode string `json:"country_code" form:"country_code" query:"country_code"`
	Mobile      string `json:"mobile" form:"mobile" query:"mobile"`
	Password    string `json:"password" form:"password" query:"password"`
}

type SignupRequest struct {
	LoginRequest
	Code string `json:"code" form:"code" query:"code"`
}

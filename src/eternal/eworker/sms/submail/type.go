package submail

const (
	StatusSuccess = "success"
)

type XSendResult struct {
	Status     string `json:"status"`
	SendID     string `json:"send_id"`
	Fee        int    `json:"fee"`
	SMSCredits string `json:"sms_credits"`
	Code       string `json:"code"`
	Msg        string `json:"msg"`
}

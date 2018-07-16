package account

const (
	KeyAccount = "account"
	KeyToken   = "token"
)

func getTokenKey(tokenID string) string {
	return KeyToken + "." + tokenID
}

func getAccountKey(userID string) string {
	return KeyAccount + "." + userID
}

package account

const (
	KeyAccount = "eternal.account"
	KeyToken   = "eternal.token"
)

func getTokenKey(tokenID string) string {
	return KeyToken + "." + tokenID
}

func getAccountKey(userID string) string {
	return KeyAccount + "." + userID
}

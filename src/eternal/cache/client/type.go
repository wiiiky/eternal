package client

const (
	KeyClient = "client"
)

func getClientKey(clientID string) string {
	return KeyClient + "." + clientID
}

package client

const (
	KeyClient = "eternal.client"
)

func getClientKey(clientID string) string {
	return KeyClient + "." + clientID
}

package opdgo

var (
	global *Client
)

func Init(client *Client) {
	global = client
}

func Global() *Client {
	return global
}

func Track(action string, properties map[string]any) {
	global.Track(action, properties)
}

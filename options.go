package opdgo

import "net/http"

type Logger interface {
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Options struct {
	Logger     Logger
	ApiURL     string
	HttpClient *http.Client
	Debug      bool
}

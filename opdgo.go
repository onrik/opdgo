package opdgo

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const (
	version        = "0.2.1"
	defaultURL     = "https://api.openpanel.dev"
	defaultTimeout = 1 * time.Second
)

type Payload struct {
	Name       string         `json:"name"`
	ProfileID  string         `json:"profile_id,omitempty"`
	Properties map[string]any `json:"properties"`
}

func (p *Payload) Add(properties map[string]any) {
	for k, v := range properties {
		p.Properties[k] = v
	}
}

type Request struct {
	Type    string  `json:"type,omitempty"`
	Payload Payload `json:"payload,omitempty"`
}

type Client struct {
	id        string
	secret    string
	baseURL   string
	logger    Logger
	http      *http.Client
	debug     bool
	global    map[string]any
	profileID string
}

func New(clientID, clientSecret string, options Options) *Client {
	if options.ApiURL == "" {
		options.ApiURL = defaultURL
	}

	if options.Logger == nil {
		options.Logger = slog.Default()
	}

	if options.HttpClient == nil {
		options.HttpClient = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	options.Logger.Debug("[opdgo] Create client", "id", clientID, "options", options)

	return &Client{
		baseURL: options.ApiURL,
		id:      clientID,
		secret:  clientSecret,
		logger:  options.Logger,
		http:    options.HttpClient,
		debug:   options.Debug,
		global:  map[string]any{},
	}
}

func (c *Client) SetGlobal(properties map[string]any) {
	for k, v := range properties {
		c.global[k] = v
	}
}

func (c *Client) ClearGlobal() {
	c.profileID = ""
	c.global = map[string]any{}
}

func (c *Client) do(requestType, action string, properties map[string]any) {
	if c == nil {
		return
	}

	payload := Payload{
		Name:       action,
		ProfileID:  c.profileID,
		Properties: map[string]any{},
	}

	payload.Add(c.global)
	payload.Add(properties)

	data, err := json.Marshal(Request{
		Type:    requestType,
		Payload: payload,
	})
	if err != nil {
		c.logger.Error("Marshal request error", "error", err)
		return
	}

	url := c.baseURL + "/track"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		c.logger.Error("Create http request error", "error", err)
		return
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("openpanel-client-id", c.id)
	request.Header.Add("openpanel-client-secret", c.secret)
	request.Header.Add("openpanel-sdk-name", "opdgo")
	request.Header.Add("openpanel-sdk-version", version)

	if c.debug {
		c.logger.Debug("[opdgo] Request->", "url", url, "body", string(data))
	}

	response, err := c.http.Do(request)
	if err != nil {
		c.logger.Error("[opdgo] Response<-", "url", url, "error", err)
		return
	}

	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)

	if response.StatusCode >= 400 {
		c.logger.Error("[opdgo] Response<-", "url", url, "status", response.StatusCode, "body", string(body))
		return
	}

	if c.debug {
		c.logger.Debug("[opdgo] Response<-", "url", url, "status", response.StatusCode, "body", string(body))
	}
}

func (c *Client) Track(action string, properties map[string]any) {
	go c.do("track", action, properties)
}

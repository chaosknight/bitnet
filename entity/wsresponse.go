package entity

import (
	"encoding/json"
	"fmt"
	"os"
)

type APIKey struct {
	ApiKey     string
	SecretKey  []byte
	Passphrase string
}

func GetApiKeyFromEnv() *APIKey {

	apikey := os.Getenv("APIKEY")
	secretkey := os.Getenv("SECRET_KEY")
	passphrase := os.Getenv("PASS_PHRASE")

	return NewAPIKey(apikey, secretkey, passphrase)
}

// NewAPIKey creates a new API key
func NewAPIKey(apiKey string, secretKey string, passphrase string) *APIKey {
	return &APIKey{
		ApiKey:     apiKey,
		SecretKey:  []byte(secretKey),
		Passphrase: passphrase,
	}
}

type Response struct {
	Event string            `json:"event,omitempty"`
	Code  JSONInt64         `json:"code"`
	Msg   string            `json:"msg"`
	ID    string            `json:"id,omitempty"`
	Op    string            `json:"op,omitempty"`
	Data  json.RawMessage   `json:"data,omitempty"`
	Arg   map[string]string `json:"arg,omitempty"`
}

func IsOpRes(m *Response) bool {

	if m.Event != "" || m.Op != "" {
		return true
	}
	return false
}

func UnmarshalResponse(b []byte) (*Response, error) {
	var m Response
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response failed :%s", err)
	}
	return &m, nil
}

func (r *Response) String() string {
	j, err := json.Marshal(r)
	if err != nil {
		return fmt.Sprintf("response to json error: %s", err)
	}
	return string(j)
}

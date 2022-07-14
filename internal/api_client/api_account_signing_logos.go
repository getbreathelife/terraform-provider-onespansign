package api_client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SigningLogo struct {
	Language string `json:"language"`
	Image    string `json:"image"`
}

const API_PATH = "/api/account/admin/signingLogos"

func (c *ApiClient) UpdateSigningLogos(d []SigningLogo) (*http.Response, error) {
	body, err := json.Marshal(d)

	if err != nil {
		return nil, err
	}

	req, err := c.newApiRequest("POST", API_PATH, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

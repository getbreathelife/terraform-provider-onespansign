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

func (c *ApiClient) UpdateSigningLogos(d []SigningLogo) *ApiError {
	var body []byte

	if len(d) > 0 {
		var err error

		body, err = json.Marshal(d)

		if err != nil {
			return &ApiError{
				Summary: "unable to marshal the request body",
				Detail:  err.Error(),
			}
		}
	} else {
		body = []byte("{}")
	}

	req, err := c.newApiRequest("POST", API_PATH, bytes.NewBuffer(body))

	if err != nil {
		return &ApiError{
			Summary: "unable to create the API request",
			Detail:  err.Error(),
		}
	}

	res, err := c.client.Do(req)

	if err != nil {
		return &ApiError{
			Summary: "unable to send the API request",
			Detail:  err.Error(),
		}
	}

	if res.StatusCode != http.StatusOK {
		return getApiError(res)
	}

	return nil
}

func (c *ApiClient) GetSigningLogos() ([]SigningLogo, *ApiError) {
	req, err := c.newApiRequest("GET", API_PATH, nil)

	if err != nil {
		return nil, &ApiError{
			Summary: "unable to create the API request",
			Detail:  err.Error(),
		}
	}

	res, err := c.client.Do(req)

	if err != nil {
		return nil, &ApiError{
			Summary: "unable to send the API request",
			Detail:  err.Error(),
		}
	}

	if res.StatusCode != http.StatusOK {
		return nil, getApiError(res)
	}

	var jsonResp []SigningLogo

	if err := jsonDecode(res.Body, &jsonResp); err != nil {
		return nil, &ApiError{
			Summary: "unable to unmarshal the API response",
			Detail:  err.Error(),
		}
	}

	return jsonResp, nil
}

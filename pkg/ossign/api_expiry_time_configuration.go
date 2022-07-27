package ossign

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ExpiryTimeConfiguration struct {
	// Default expiry time for transactions in days
	Default json.Number `json:"remainingDays"`

	// Maximum allowed value for expiry time for transactions in days
	Maximum json.Number `json:"maximumRemainingDays"`
}

func (c *ApiClient) GetExpiryTimeConfiguration() (*ExpiryTimeConfiguration, *ApiError) {
	res, err := c.makeApiRequest("GET", "/api/dataRetentionSettings/expiryTimeConfiguration", nil)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, getApiError(res)
	}

	var jsonResp ExpiryTimeConfiguration

	if err := jsonDecode(res.Body, &jsonResp); err != nil {
		return nil, &ApiError{
			Summary: "unable to unmarshal the API response",
			Detail:  err.Error(),
		}
	}

	return &jsonResp, nil
}

func (c *ApiClient) UpdateExpiryTimeConfiguration(d ExpiryTimeConfiguration) *ApiError {
	body, err := json.Marshal(d)

	if err != nil {
		return &ApiError{
			Summary: "unable to marshal the request body",
			Detail:  err.Error(),
		}
	}

	res, apiErr := c.makeApiRequest("PUT", "/api/dataRetentionSettings/expiryTimeConfiguration", bytes.NewBuffer(body))

	if apiErr != nil {
		return apiErr
	}

	if res.StatusCode != http.StatusOK {
		return getApiError(res)
	}

	return nil
}

package ossign

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type TransactionRetention struct {
	// Number of days to keep drafts for
	Draft int64 `json:"draft"`

	// Number of days to keep sent transactions for
	Sent int64 `json:"sent"`

	// Number of days to keep completed transactions for
	Completed int64 `json:"completed"`

	// Number of days to keep archived transactions for
	Archived int64 `json:"archived"`

	// Number of days to keep declined transactions for
	Declined int64 `json:"declined"`

	// Number of days to keep opted-out transactions for
	OptedOut int64 `json:"optedOut"`

	// Number of days to keep expired transactions for
	Expired int64 `json:"expired"`
}

type DataManagementPolicy struct {
	TransactionRetention TransactionRetention `json:"transactionRetention"`
}

type ExpiryTimeConfiguration struct {
	// Default expiry time for transactions in days
	Default int64 `json:"remainingDays"`

	// Maximum allowed value for expiry time for transactions in days
	Maximum int64 `json:"maximumRemainingDays"`
}

func (c *ApiClient) GetDataManagementPolicy() (*DataManagementPolicy, *ApiError) {
	res, err := c.makeApiRequest("GET", "/api/dataRetentionSettings/dataManagementPolicy", nil)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, getApiError(res)
	}

	var jsonResp DataManagementPolicy

	if err := jsonDecode(res.Body, &jsonResp); err != nil {
		return nil, &ApiError{
			Summary: "unable to unmarshal the API response",
			Detail:  err.Error(),
		}
	}

	return &jsonResp, nil
}

func (c *ApiClient) UpdateDataManagementPolicy(d DataManagementPolicy) *ApiError {
	body, err := json.Marshal(d)

	if err != nil {
		return &ApiError{
			Summary: "unable to marshal the request body",
			Detail:  err.Error(),
		}
	}

	res, apiErr := c.makeApiRequest("PUT", "/api/dataRetentionSettings/dataManagementPolicy", bytes.NewBuffer(body))

	if apiErr != nil {
		return apiErr
	}

	if res.StatusCode != http.StatusOK {
		return getApiError(res)
	}

	return nil
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

	res, apiErr := c.makeApiRequest("PUT", "/api/dataRetentionSettings/dataManagementPolicy", bytes.NewBuffer(body))

	if apiErr != nil {
		return apiErr
	}

	if res.StatusCode != http.StatusOK {
		return getApiError(res)
	}

	return nil
}

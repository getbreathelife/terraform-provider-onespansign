package ossign

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type TransactionRetention struct {
	// Number of days to keep drafts for
	Draft json.Number `json:"draft"`

	// Number of days to keep sent transactions for
	Sent json.Number `json:"sent"`

	// Number of days to keep completed transactions for
	Completed json.Number `json:"completed"`

	// Number of days to keep archived transactions for
	Archived json.Number `json:"archived"`

	// Number of days to keep declined transactions for
	Declined json.Number `json:"declined"`

	// Number of days to keep opted-out transactions for
	OptedOut json.Number `json:"optedOut"`

	// Number of days to keep expired transactions for
	Expired json.Number `json:"expired"`

	// Number of days to keep the transactions, calculated from the day that the
	// transaction is created
	LifetimeTotal json.Number `json:"lifetimeTotal"`

	// Number of days that incomplete transactions will be stored, calculated from
	// the day that the transaction is created
	LifetimeUntilCompletion json.Number `json:"lifetimeUntilCompletion"`

	// Include sent transactions as part of the "incomplete transactions"
	IncludeSent bool `json:"includeSent"`
}

type DataManagementPolicy struct {
	TransactionRetention TransactionRetention `json:"transactionRetention"`
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

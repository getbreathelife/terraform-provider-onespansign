package ossign

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SigningTheme struct {
	// Primary color hex value
	Primary string `json:"primary"`

	// Success color hex value
	Success string `json:"success"`

	// Warning color hex value
	Warning string `json:"warning"`

	// Error color hex value
	Error string `json:"error"`

	// Info color hex value
	Info string `json:"info"`

	// Required signature button color hex value
	SignatureButton string `json:"signatureButton"`

	// Optional signature button color hex value
	OptionalSignatureButton string `json:"optionalSignatureButton"`
}

// CreateAccountSigningThemes creates customized signing themes on the account.
//
// https://community.onespan.com/products/onespan-sign/sandbox#/Account%20Signing%20Themes/api.account.signingThemes.post
func (c *ApiClient) CreateAccountSigningThemes(t map[string]SigningTheme) *ApiError {
	body, err := json.Marshal(t)

	if err != nil {
		return &ApiError{
			Summary: "unable to marshal the request body",
			Detail:  err.Error(),
		}
	}

	res, apiErr := c.makeApiRequest("POST", "/api/account/signingThemes", bytes.NewBuffer(body))

	if apiErr != nil {
		return apiErr
	}

	if res.StatusCode != http.StatusOK {
		return getApiError(res)
	}

	return nil
}

// GetAccountSigningThemes retrieves the customized signing themes on the account.
//
// https://community.onespan.com/products/onespan-sign/sandbox#/Account%20Signing%20Themes/api.account.signingThemes.get
func (c *ApiClient) GetAccountSigningThemes() (map[string]SigningTheme, *ApiError) {
	res, err := c.makeApiRequest("GET", "/api/account/signingThemes", nil)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, getApiError(res)
	}

	var jsonResp map[string]SigningTheme

	if err := jsonDecode(res.Body, &jsonResp); err != nil {
		return nil, &ApiError{
			Summary: "unable to unmarshal the API response",
			Detail:  err.Error(),
		}
	}

	return jsonResp, nil
}

// UpdateAccountSigningThemes updates the customized signing themes on the account.
//
// https://community.onespan.com/products/onespan-sign/sandbox#/Account%20Signing%20Themes/api.account.signingThemes.put
func (c *ApiClient) UpdateAccountSigningThemes(t map[string]SigningTheme) *ApiError {
	body, err := json.Marshal(t)

	if err != nil {
		return &ApiError{
			Summary: "unable to marshal the request body",
			Detail:  err.Error(),
		}
	}

	res, apiErr := c.makeApiRequest("PUT", "/api/account/signingThemes", bytes.NewBuffer(body))

	if apiErr != nil {
		return apiErr
	}

	if res.StatusCode != http.StatusOK {
		return getApiError(res)
	}

	return nil
}

// DeleteAccountSigningThemes deletes the customized signing themes on the account.
//
// https://community.onespan.com/products/onespan-sign/sandbox#/Account%20Signing%20Themes/api.account.signingThemes.put
func (c *ApiClient) DeleteAccountSigningThemes() *ApiError {
	res, apiErr := c.makeApiRequest("DELETE", "/api/account/signingThemes", nil)

	if apiErr != nil {
		return apiErr
	}

	if res.StatusCode != http.StatusOK {
		return getApiError(res)
	}

	return nil
}

func (l SigningTheme) Equal(r SigningTheme) bool {
	return l.Primary == r.Primary &&
		l.Success == r.Success &&
		l.Warning == r.Warning &&
		l.Error == r.Error &&
		l.Info == r.Info &&
		l.SignatureButton == r.SignatureButton &&
		l.OptionalSignatureButton == r.OptionalSignatureButton
}

package ossign

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SigningLogo struct {
	Language string `json:"language"`
	Image    string `json:"image"`
}

// UpdateAccountSigningLogos Adds, updates or deletes an account's customized Signing Ceremony logos.
//
// https://community.onespan.com/products/onespan-sign/sandbox#/Account%20Signing%20Logos/api.account.admin.signingLogos.post
func (c *ApiClient) UpdateAccountSigningLogos(d []SigningLogo) *ApiError {
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
		body = []byte("[]")
	}

	res, err := c.makeApiRequest("POST", "/api/account/admin/signingLogos", bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return getApiError(res)
	}

	return nil
}

// GetAccountSigningLogos Retrieves an account's customized logo for use during the Signing Ceremony. In addition, the corresponding langauge for the account is also retrieved.
//
// https://community.onespan.com/products/onespan-sign/sandbox#/Account%20Signing%20Logos/api.account.admin.signingLogos.get
func (c *ApiClient) GetAccountSigningLogos() ([]SigningLogo, *ApiError) {
	res, err := c.makeApiRequest("GET", "/api/account/admin/signingLogos", nil)

	if err != nil {
		return nil, err
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

func (l SigningLogo) Equal(r SigningLogo) bool {
	return l.Language == r.Language && l.Image == r.Image
}

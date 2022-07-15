package ossign

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type ApiClientConfig struct {
	BaseUrl      string
	ClientId     string
	ClientSecret string
	UserAgent    string
}

type ApiClient struct {
	baseUrl string
	ua      string
	token   string
	expiry  int64

	client *http.Client

	ClientId     string
	clientSecret string
}

type ApiError struct {
	Summary string
	Detail  string
}

type accessTokenResponse struct {
	AccessToken string      `json:"accessToken"`
	ExpiresAt   json.Number `json:"expiresAt"`
}

type ErrorResponse struct {
	Code       json.Number `json:"code"`
	MessageKey string      `json:"messageKey"`
	Message    string      `json:"message"`
}

const API_VERSION = "11.47"

func NewClient(config *ApiClientConfig) *ApiClient {
	client := &http.Client{}

	return &ApiClient{
		baseUrl:      config.BaseUrl,
		ua:           config.UserAgent,
		client:       client,
		ClientId:     config.ClientId,
		clientSecret: config.ClientSecret,
	}
}

func (c *ApiClient) getAuthToken() (string, error) {
	if c.token != "" && c.expiry > time.Now().Unix() {
		return c.token, nil
	}

	url := &url.URL{
		Scheme: "https",
		Host:   c.baseUrl,
		Path:   "/apitoken/clientApp/accessToken",
	}

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer([]byte(fmt.Sprintf(`{"clientId":"%s","secret":"%s","type":"OWNER"}`, c.ClientId, c.clientSecret))))
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.ua)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}

	var jsonResp accessTokenResponse

	if err := jsonDecode(resp.Body, &jsonResp); err != nil {
		return "", err
	}

	c.token = jsonResp.AccessToken
	c.expiry, err = jsonResp.ExpiresAt.Int64()

	if err != nil {
		return "", err
	}

	if c.token == "" {
		return "", errors.New("unable to retrieve an access token for OneSpan's API")
	}

	return c.token, nil
}

// newApiRequest creates a HTTP request to the OneSpan Sign API host configured in the ApiClient.
// It accepts a HTTP method string, path (not full URL) to the API resource, and the request body.
// It also automatically retrieves the access token for the API and inserts it to the request's Authorization header.
func (c *ApiClient) newApiRequest(method string, path string, body io.Reader) (*http.Request, error) {
	token, err := c.getAuthToken()
	if err != nil {
		return nil, err
	}

	url := &url.URL{
		Scheme: "https",
		Host:   c.baseUrl,
		Path:   path,
	}

	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", fmt.Sprintf("application/json; esl-api-version=%s", API_VERSION))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.ua)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	return req, nil
}

func UnmarshalApiErrorResponse(res *http.Response) (*ErrorResponse, error) {
	var jsonResp ErrorResponse

	if err := jsonDecode(res.Body, &jsonResp); err != nil {
		return nil, err
	}

	return &jsonResp, nil
}

// jsonDecode creates a JSON decoder that reads from r and stores the decoded value in the value pointed to by v
func jsonDecode(r io.Reader, v interface{}) error {
	d := json.NewDecoder(r)
	d.UseNumber()
	return d.Decode(v)
}

func getApiError(res *http.Response) *ApiError {
	apiErr := &ApiError{
		Summary: "invalid response from the API",
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		apiErr.Detail = "unable to read the response"
		return apiErr
	}

	var detail bytes.Buffer

	if err = json.Indent(&detail, body, "", "\t"); err == nil {
		apiErr.Detail = "unable to parse the response"
		return apiErr
	}

	errMsg, err := io.ReadAll(&detail)
	if err != nil {
		apiErr.Detail = "unable to read the parsed response"
		return apiErr
	}

	apiErr.Detail = string(errMsg)
	return apiErr
}
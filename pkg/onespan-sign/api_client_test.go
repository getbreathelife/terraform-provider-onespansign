package ossign_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	ossign "github.com/getbreathelife/terraform-provider-onespan-sign/pkg/onespan-sign"
	"github.com/getbreathelife/terraform-provider-onespan-sign/pkg/onespan-sign/testhelpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type testServerConfig struct {
	// Mock access token return from the api token endpoint
	AccessToken string
	// TokenExpiryOffset is the amount of seconds that should be given until the api token is expired
	TokenExpiryOffset int64
}

const testImg = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAoAAAANCAYAAACQN/8FAAAABGdBTUEAALGPC" +
	"/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAhGVYSWZNTQAqAAAACAAFARIAAwAAAAEA" +
	"AQAAARoABQAAAAEAAABKARsABQAAAAEAAABSASgAAwAAAAEAAgAAh2kABAAAAAEAAABaAAAAAAAAASAAAAABAAABIAAAAAEAA6A" +
	"BAAMAAAABAAEAAKACAAQAAAABAAAACqADAAQAAAABAAAADQAAAADcXFczAAAACXBIWXMAACxLAAAsSwGlPZapAAABWWlUWHRYTU" +
	"w6Y29tLmFkb2JlLnhtcAAAAAAAPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iWE1QIENvcmUgN" +
	"i4wLjAiPgogICA8cmRmOlJERiB4bWxuczpyZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMi" +
	"PgogICAgICA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIgogICAgICAgICAgICB4bWxuczp0aWZmPSJodHRwOi8vbnMuYWR" +
	"vYmUuY29tL3RpZmYvMS4wLyI+CiAgICAgICAgIDx0aWZmOk9yaWVudGF0aW9uPjE8L3RpZmY6T3JpZW50YXRpb24+CiAgICAgID" +
	"wvcmRmOkRlc2NyaXB0aW9uPgogICA8L3JkZjpSREY+CjwveDp4bXBtZXRhPgoZXuEHAAAByElEQVQoFS1Rz2sTQRR+b2Z2N7Ek2" +
	"mpyyEHi1RQ8VA+lBXuyzcGbK3gTLxZJpf+A0IOIV7HtoVcRweCpWBoUBCWe7HEREVpFLNg1lCbZbfbHzPNNcGB4v775+N43uNjq" +
	"joR0PakKYIwGncWvyhfO3m2vNVIgQkAk4KNQOJ7J019G5+8BqFYs1+70//4+4tnqzP09tQeQWaBQThFI4E5nffYe3+Zp/7DP/aY" +
	"d1gAc33/t2lxlySBHgHlbACAZ070OSka22t66Gtvo+yRxaaX7UDoTz/Ik3gWEb6wp4ofGEEyxlBPH4Iu3m7Nfcan1aUG6pQ9Cen" +
	"YRYM2WBIxOwS2eg9HwiKXBnCCQz/M01kncu0Eaq9pElygZ1MHE1dHwzy3HKwMaeqqcQnk6Pz152dmYfzemAgj/RxveLLY+f2Tt0" +
	"4rYO97a2K6/FrgQBDpsVLAShNRu39ZsIpsJmUIUDKHEAkeH+2rQqJgKVMX+ZMkarS2IzZAqS4dc4xUL3N66ObbD5vY0V757hsLL" +
	"QBAr/qDNM5P1B6zlC88CbhaYgdkw0Sa8Vjpfr0W9g0fqYuqs/jz+IVjCsnInZojGchmHoJPIDHsHTzobc4//ASnLxOSzgBDYAAA" +
	"AAElFTkSuQmCC"

// setupTestServer sets up a server that emulates a OneSpan Sign server.
// Returns a http.Request pointer that corresponds to the last received request, and the pointer to the test server.
func setupTestServer(tsc *testServerConfig) (*testhelpers.HttpRequestHistory, *httptest.Server) {
	h := testhelpers.NewHttpRequestHistory()

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}

			h.Push(&testhelpers.RequestHistoryEntry{
				Request: r,
				Body:    b,
			})

			switch r.URL.Path {
			case "/api/account/admin/signingLogos":
				switch r.Method {
				case "GET":
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode([]map[string]interface{}{
						{
							"language": "en",
							"image":    testImg,
						},
						{
							"language": "fr",
							"image":    testImg,
						},
					})

				default:
					w.WriteHeader(http.StatusNotFound)
				}

			case "/apitoken/clientApp/accessToken":
				switch r.Method {
				case "POST":
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"accessToken": tsc.AccessToken,
						"expiresAt":   time.Now().Unix() + tsc.TokenExpiryOffset,
					})

				default:
					w.WriteHeader(http.StatusNotFound)
				}
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}),
	)

	return h, ts
}

func TestGetAuthToken(t *testing.T) {
	token := uuid.NewString()

	h, ts := setupTestServer(&testServerConfig{
		AccessToken:       token,
		TokenExpiryOffset: 5,
	})
	defer ts.Close()

	url, err := url.Parse(ts.URL)

	if err != nil {
		panic(err)
	}

	cid := uuid.NewString()
	csct := uuid.NewString()

	c := ossign.NewClient(ossign.ApiClientConfig{
		BaseUrl:      url,
		ClientId:     cid,
		ClientSecret: csct,
		UserAgent:    uuid.NewString(),
	})

	c.GetSigningLogos()
	c.GetSigningLogos()

	assert.Equal(t, 3, len(h.Stack))

	r := h.Stack[0]

	assert.Equal(t, "/apitoken/clientApp/accessToken", r.Request.URL.Path)
	assert.Equal(t, "application/json", r.Request.Header.Get("Accept"))
	assert.Equal(t, "application/json", r.Request.Header.Get("Content-Type"))
	assert.Empty(t, r.Request.Header.Get("Authorization"))

	assert.JSONEq(t, fmt.Sprintf(`{"clientId":"%s","secret":"%s","type":"OWNER"}`, cid, csct), string(r.Body))

	r = h.Stack[1]
	assert.Equal(t, fmt.Sprintf("Bearer %s", token), r.Request.Header.Get("Authorization"))

	r = h.Stack[2]
	assert.Equal(t, fmt.Sprintf("Bearer %s", token), r.Request.Header.Get("Authorization"))
}

func TestGetAuthTokenRefresh(t *testing.T) {
	token := uuid.NewString()

	config := &testServerConfig{
		AccessToken:       token,
		TokenExpiryOffset: 0,
	}

	h, ts := setupTestServer(config)
	defer ts.Close()

	url, err := url.Parse(ts.URL)

	if err != nil {
		panic(err)
	}

	cid := uuid.NewString()
	csct := uuid.NewString()

	c := ossign.NewClient(ossign.ApiClientConfig{
		BaseUrl:      url,
		ClientId:     cid,
		ClientSecret: csct,
		UserAgent:    uuid.NewString(),
	})

	c.GetSigningLogos()

	token2 := uuid.NewString()
	config.AccessToken = token2

	// Wait until the token is expired
	time.Sleep(time.Second)

	c.GetSigningLogos()

	assert.Equal(t, 4, len(h.Stack))

	r := h.Stack[0]

	assert.Equal(t, "/apitoken/clientApp/accessToken", r.Request.URL.Path)
	assert.Equal(t, "application/json", r.Request.Header.Get("Accept"))
	assert.Equal(t, "application/json", r.Request.Header.Get("Content-Type"))
	assert.Empty(t, r.Request.Header.Get("Authorization"))

	assert.JSONEq(t, fmt.Sprintf(`{"clientId":"%s","secret":"%s","type":"OWNER"}`, cid, csct), string(r.Body))

	r = h.Stack[1]
	assert.Equal(t, fmt.Sprintf("Bearer %s", token), r.Request.Header.Get("Authorization"))

	r = h.Stack[2]
	assert.Equal(t, "/apitoken/clientApp/accessToken", r.Request.URL.Path)

	r = h.Stack[3]
	assert.Equal(t, fmt.Sprintf("Bearer %s", token2), r.Request.Header.Get("Authorization"))
}

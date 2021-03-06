package insta

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	// Base Instagram API URL
	base = "https://api.instagram.com/v1"
)

// InstaCLient gives access to the Instagram API
type InstaClient struct {
	AccessToken string
	Client      *http.Client
	ClientID    string
}

// NewInstaClient returns an initialized InstaClient, that uses the
// given access token and http.Client to make requests
func NewClient(client *http.Client, accessToken string) *InstaClient {
	c := new(InstaClient)
	c.AccessToken = accessToken
	c.Client = client
	return c
}

// getRequest sends a GET request to the Instagram API and unmarshals the received data.
func (i *InstaClient) getRequest(endpoint string, options map[string]string, resultType interface{}) error {
	// Convert the options into URL query string
	urlParameters := url.Values{}
	urlParameters.Add("access_token", i.AccessToken)
	for key, value := range options {
		urlParameters.Add(key, value)
	}

	// Convert full request url into URL struct, so we can add the query string.
	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	// Add query string to request url
	u.RawQuery = urlParameters.Encode()

	// Send request
	resp, err := i.Client.Get(u.String())
	if err != nil {
		return err
	}
	// Check response code
	if resp.StatusCode != 200 {
		return newAPIError(resp)
	}

	// Decode JSON response into given struct
	err = decodeBody(resp.Body, resultType)
	if err != nil {
		return err
	}

	return nil
}

// APIError represents an error originating from the Instagram API.
type APIError ResponseMeta

// Error returns the string representation of the error.
func (a APIError) Error() string {
	return fmt.Sprintf("Instagram API error: Code: %d, Type: %s, Message: %s", a.Code, a.ErrorType, a.ErrorMessage)
}

// newAPIError returns the error sent by the Instagram API.
func newAPIError(r *http.Response) APIError {
	var meta ApiResponse
	decodeBody(r.Body, &meta)

	return APIError(meta.Meta)
}

// decodeBody decoodes a HTTP response's body into the requested type
func decodeBody(body io.Reader, resultType interface{}) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(resultType)
	if err != nil {
		return err
	}

	return nil
}

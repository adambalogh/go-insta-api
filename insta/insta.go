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
	instagramApiBaseURL = "https://api.instagram.com/v1"
)

// InstaCLient gives access to the Instagram API client
// It normally it requires an access token, but some parts of the API can be
// accessed by just using the client ID, please check the Instagram API doc.
type InstaClient struct {
	client   *http.Client
	ClientID string
}

// NewInstaClient returns an initialized InstaClient, with a built-in HTTPClient
func NewClient(client *http.Client) *InstaClient {
	c := new(InstaClient)
	c.client = client
	return c
}

// getRequest dispatches a GET request to the Instagram API and unmarshals the received data.
func (i *InstaClient) getRequest(endpoint string, options map[string]string, resultType interface{}) error {
	// Convert the options into URL query string
	urlParameters := url.Values{}
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
	resp, err := i.client.Get(u.String())
	if err != nil {
		return err
	}
	// Check response code
	if resp.StatusCode != 200 {
		return newApiError(resp)
	}

	// Decode JSON response into given struct
	err = decodeBody(resp.Body, resultType)
	if err != nil {
		return err
	}

	return nil
}

// ApiError represents an error originating from the Instagram API.
type ApiError ResponseMeta

// Error returns the string representation of the error.
func (a ApiError) Error() string {
	return fmt.Sprintf("Instagram API error: Code: %d, Type: %s, Message: %s", a.Code, a.ErrorType, a.ErrorMessage)
}

// newApiError returns the error sent by the Instagram API.
func newApiError(r *http.Response) ApiError {
	var meta ApiResponse
	decodeBody(r.Body, &meta)

	return ApiError(meta.Meta)
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

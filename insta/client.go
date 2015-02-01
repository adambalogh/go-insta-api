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

// HTTPRequester sends HTTP requests
type HTTPRequester interface {
	SendGetRequest(url string) (*http.Response, error)
	SendPostRequest(url string, formValues url.Values) (*http.Response, error)
}

// SimpleHTTPRequester sends HTTP requests using the built-in library
type SimpleHTTPRequester struct {
}

// SendRequest sends a HTTP GET request to the requested URL
func (s SimpleHTTPRequester) SendGetRequest(url string) (*http.Response, error) {
	// Send request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SendPostRequest sends a HTTP POST request with the given form values
// to the requested URL
func (s SimpleHTTPRequester) SendPostRequest(url string, formValues url.Values) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.PostForm(url, formValues)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Instagram API client, it normally it requires an access
// token, but some parts of the API can be accessed by just
// using the client ID, please check the Instagram API doc.
type InstaClient struct {
	HTTPRequester
	ClientID    string
	AccessToken string
}

// NewInstaClient returns an initialized InstaClient, with
// a basic HTTPRequester
func NewInstaClient(accessToken string) *InstaClient {
	client := new(InstaClient)
	client.HTTPRequester = SimpleHTTPRequester{}
	client.AccessToken = accessToken
	return client
}

// getRequest sends GET request to Instagram API and unmarshals received data.
func (i *InstaClient) getRequest(endpointURL string, options map[string]string, resultType interface{}) error {
	// Convert the options into URL query string
	urlParameters := url.Values{}
	for key, value := range options {
		urlParameters.Add(key, value)
	}
	// TODO not all endpoints require access tokens
	urlParameters.Add("access_token", i.AccessToken)

	// Convert full request url into URL struct, so we can add the query string.
	completeURL := instagramApiBaseURL + endpointURL
	u, err := url.Parse(completeURL)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Add query string to request url
	u.RawQuery = urlParameters.Encode()
	//fmt.Println(u.String())

	// Send request
	resp, err := i.SendGetRequest(u.String())
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
func newApiError(r *http.Response) error {
	var meta ApiResponse
	decodeBody(r.Body, &meta)

	return ApiError(meta.Meta)
}

// decodeBody decoodes a HTTP response's body into the requested interface type
func decodeBody(body io.Reader, resultType interface{}) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(resultType)
	if err != nil {
		return err
	}

	return nil
}

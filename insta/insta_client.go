package insta

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// Base Instagram API URL
	instagramApiBaseUrl = "https://api.instagram.com/v1"
)

// Instagram API client, it normally it requires an access
// token, but some parts of the API can be accessed by just
// using the client ID, please check the Instagram API doc
type InstaClient struct {
	ClientId    string
	AccessToken string
}

// Sends request to Instagram API and unmarshals received data
func (i *InstaClient) requestUrl(endpointUrl string, options map[string]string, resultType interface{}) error {
	// Convert the options into URL values
	urlParameters := url.Values{}
	// TODO not all endpoints require access tokens
	urlParameters.Add("access_token", i.AccessToken)
	for key, value := range options {
		urlParameters.Add(key, value)
	}

	// Convert full API url into URL struct, so we can
	// add the query string
	completeUrl := instagramApiBaseUrl + endpointUrl
	u, err := url.Parse(completeUrl)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Add query string to url
	u.RawQuery = urlParameters.Encode()
	fmt.Println(u.String())

	// Send request to Instagram API
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	// Check response code
	if resp.StatusCode != 200 {
		// TODO Extract error message from API response
		return newApiError(resp)
	}
	
	// Unmarshal JSON into given struct type
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&resultType)
	if err != nil {
		return err
	}
	
	return nil
}

// Represents an error originating from the Instagram API
type ApiError ResponseMeta

func (a ApiError) Error() string {
	return fmt.Sprintf("Instagram API error: Code: %d, Type: %s, Message: %s", a.Code, a.ErrorType, a.ErrorMessage)
}

// Returns the error sent by the Instagram APi
func newApiError(r *http.Response) error {
	var meta ApiResponse
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&meta)	
	if err != nil {
		return err
	}
	
	return ApiError(meta.Meta)
}

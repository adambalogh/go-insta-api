package insta

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// Base Instagram API URL
	instagramApiBaseUrl = "https://api.instagram.com/v1/"
)

// Instagram API client, it normally it requires an access
// token, but some parts of the API can be accessed by just
// using the client ID, please check the Instagram API doc
type InstaClient struct {
	ClientId    string
	AccessToken string
}

// Send request to Instagram API and unmarshal received data
func (i *InstaClient) get(endpointUrl string, options map[string]string, resultType interface{}) error {
	// Encode options to Query String
	urlParameters := url.Values{}
	// TODO not all encdpoints require access tokens
	urlParameters.Add("access_token", i.AccessToken)
	for key, value := range options {
		urlParameters.Add(key, value)
	}

	// Parse Instagram API url
	completeUrl = instagramApiBaseUrl + endpointUrl
	u, err := url.Parse(completeUrl)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Add query string to url
	u.RawQuery = urlParameters.Encode()

	// Get posts from given URL
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	// Decode JSON to get posts
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&resultType)
	if err != nil {
		return err
	}
	return nil
}

package insta

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// Instagram login page URL
	authorizationURL = "https://api.instagram.com/oauth/authorize/?client_id=%s&redirect_uri=%s&response_type=code"
	// Exchange code for access token URL
	accessTokenURL = "https://api.instagram.com/oauth/access_token"
)

// This struct is solely used to retrieve access
// token to authenticate InstaClient class
type InstaLogin struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// GetLoginURL returns Instagram login page's URL
func (n *InstaLogin) GetLoginURL() string {
	return fmt.Sprintf(authorizationURL, n.ClientID, n.RedirectURL)
}

// ExchangeCodeForAccessToken exchanges the code received from Instagram's Login page
// for an access token
func (n *InstaLogin) ExchangeCodeForAccessToken(code string) (string, error) {
	// HTTP POST form values required for code exchange
	params := url.Values{}
	params.Add("code", code)
	params.Add("client_id", n.ClientID)
	params.Add("client_secret", n.ClientSecret)
	params.Add("grant_type", "authorization_code")
	params.Add("redirect_uri", n.RedirectURL)
	// Send HTTP POST request
	client := &http.Client{}
	resp, err := client.PostForm(accessTokenURL, params)
	if err != nil {
		return "", err
	}
	// Check status code
	if resp.StatusCode != 200 {
		return "", errors.New("Failed to authenticate")
	}

	// Decode JSON to get AccessToken
	decoder := json.NewDecoder(resp.Body)
	var accessToken AccessTokenResponse
	err = decoder.Decode(&accessToken)
	if err != nil {
		return "", err
	}

	return accessToken.AccessToken, nil
}

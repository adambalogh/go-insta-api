package insta

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// Instagram login page URL
	authorizationUrl = "https://api.instagram.com/oauth/authorize/?client_id=%s&redirect_uri=%s&response_type=code"
	// Exchange code for access token URL
	accessTokenUrl = "https://api.instagram.com/oauth/access_token"
)

// Instagram access token response
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// This class is solely used to retrieve access
// token to authenticate InstaClient class
type InstaLogin struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}

// Return Instagram Login page's URL
func (n *InstaLogin) GetLoginUrl() string {
	return fmt.Sprintf(authorizationUrl, n.ClientId, n.RedirectUrl)
}

// Get Access token using code returned from Instagram's
// login page
func (n *InstaLogin) ExchangeCodeForAccessToken(code string) (string, error) {
	// HTTP POST form values required for code exchange
	params := url.Values{}
	params.Add("code", code)
	params.Add("client_id", n.ClientId)
	params.Add("client_secret", n.ClientSecret)
	params.Add("grant_type", "authorization_code")
	params.Add("redirect_uri", n.RedirectUrl)
	// Send HTTP POST request
	client := &http.Client{}
	resp, err := client.PostForm(accessTokenUrl, params)
	if err != nil {
		return "", err
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

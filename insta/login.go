package insta

import (
	"net/http"

	"golang.org/x/oauth2"
)

const (
	// Instagram login page URL
	authorizationURL = "https://api.instagram.com/oauth/authorize"
	// Exchange code for access token URL
	accessTokenURL = "https://api.instagram.com/oauth/access_token"
)

// This struct is solely used to retrieve access
// token to authenticate InstaClient class
type InstaLogin struct {
	HTTPClient   *http.Client
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Conf         *oauth2.Config
	Context      *oauth2.Context
}

// New returns an initialized InstaLogin, with a SimpleHTTPRequester
func New(clientID, clientSecret, redirectURL string) *InstaLogin {
	login := new(InstaLogin)
	login.HTTPClient = &http.Client{}
	login.ClientID = clientID
	login.ClientSecret = clientSecret
	login.RedirectURL = redirectURL
	login.Conf = &oauth2.Config{
		ClientID:     login.ClientID,
		ClientSecret: login.ClientSecret,
		Scopes:       []string{"basic"},
		RedirectURL:  login.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authorizationURL,
			TokenURL: accessTokenURL,
		},
	}
	// TODO allow custom context
	login.Context = &oauth2.NoContext
	return login
}

// GetLoginURL returns Instagram login page's URL
func (i *InstaLogin) GetLoginURL() string {
	return i.Conf.AuthCodeURL("test", oauth2.AccessTypeOnline)
}

// ExchangeCodeForAccessToken exchanges the code received from Instagram's Login page
// for an access token
func (i *InstaLogin) Exchange(code string) (string, error) {
	accessToken, err := i.Conf.Exchange(i.Context, code)
	if err != nil {
		return "", err
	}
	return accessToken.AccessToken, nil
}

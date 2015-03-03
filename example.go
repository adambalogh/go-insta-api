package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/adambalogh/go-insta-api/insta"
	"golang.org/x/oauth2"
)

const (
	clientID         = "2c28931a275a410aa440fa062c1a6a9d"
	clientSecret     = "c9b74f19b20d4bed80773b9f7e06696d"
	redirectURL      = "http://localhost:8080/oauth-complete"
	authorizationURL = "https://api.instagram.com/oauth/authorize"
	accessTokenURL   = "https://api.instagram.com/oauth/access_token"
)

var (
	auth = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"basic"},
		RedirectURL:  "http://localhost:8080/oauth-complete",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authorizationURL,
			TokenURL: accessTokenURL,
		},
	}
)

/*
  Search for users
*/
func SearchUser(w http.ResponseWriter, r *http.Request) {
	token, _ := (CookieTokenSource{r}).Token()
	client := insta.NewClient(&http.Client{}, token.AccessToken)
	queryString := r.URL.Query().Get("q")
	// Make request to Instagram API
	searchResult, err := client.SearchUser(queryString, map[string]string{})
	if err != nil {
		fmt.Fprintf(w, "Sorry an error occured: %s", err)
		return
	}

	// Create HTML page
	htmlStart := "<html><head></head><body>"
	htmlEnd := "</body></html>"
	htmlBody := bytes.Buffer{}
	for _, user := range searchResult.Users {
		htmlBody.WriteString("<p><a href=\"http://www.instagram.com/" +
			user.Username + "\" />" +
			user.Username + "</a></p>")
	}
	// Return HTML page
	fmt.Fprintf(w, htmlStart+htmlBody.String()+htmlEnd)
}

func UserGroup(w http.ResponseWriter, r *http.Request) {
	token, _ := (CookieTokenSource{r}).Token()
	client := insta.NewClient(&http.Client{}, token.AccessToken)
	usernames := strings.Split(r.URL.Query().Get("names"), " ")

	var userIDs []string
	for _, username := range usernames {
		userID, _ := client.GetUserID(username)
		userIDs = append(userIDs, userID)
	}

	posts, err := client.GetPostsFromUsers(userIDs, map[string]string{
		"count": "5",
	})
	if err != nil {
		fmt.Println(err)
	}

	// Create HTML page
	htmlStart := "<html><head></head><body>"
	htmlEnd := "</body></html>"
	htmlBody := bytes.Buffer{}
	for _, post := range posts {
		htmlBody.WriteString("<p><img src=\"" +
			post.Images.StandardResolution.URL + "\"/></p>")
	}
	// Return HTML page
	fmt.Fprintf(w, htmlStart+htmlBody.String()+htmlEnd)

}

/*
 Redirect user to Instagram's login page
*/
func RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	url := auth.AuthCodeURL("lolololo", oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusFound)
}

/*
 Exchange code for access token
*/
func GetAccessToken(w http.ResponseWriter, r *http.Request) {
	// Get user code, sent by Instagram Login page
	urlParameters := r.URL.Query()
	code := urlParameters.Get("code")

	token, err := auth.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Print(w, err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "a",
		Value: token.AccessToken,
	})

	fmt.Fprint(w, "done")
}

type CookieTokenSource struct {
	r *http.Request
}

func (c CookieTokenSource) Token() (*oauth2.Token, error) {
	a, _ := c.r.Cookie("a")

	token := oauth2.Token{
		AccessToken: a.Value,
		TokenType:   "Bearer",
	}
	return &token, nil
}

func main() {
	// Handle authentication
	http.HandleFunc("/login", RedirectToLogin)
	http.HandleFunc("/oauth-complete", GetAccessToken)

	http.HandleFunc("/search", SearchUser)
	http.HandleFunc("/group", UserGroup)

	http.ListenAndServe(":8080", nil)
}

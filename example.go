package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/adambalogh/go-insta-api/insta"
)

const (
	clientID     = ""
	clientSecret = ""
	redirectURL  = ""
)

// authApi is used for logging in the user, authenticating
// the application and getting the access token
var authApi insta.InstaLogin

/*
  Search for users
*/
func SearchUser(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromCookie(r)
	// Create authenticated Instagram CLient
	client := insta.NewInstaClient(token)

	queryString := r.URL.Query().Get("q")

	var searchResult *insta.SearchResult // SearchUser returns a pointer to a SearchResult struct
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

/*
 Redirect user to Instagram's login page
*/
func RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	url := authApi.GetLoginURL()
	http.Redirect(w, r, url, http.StatusFound)
}

/*
 Exchange code for access token
*/
func GetAccessToken(w http.ResponseWriter, r *http.Request) {
	// Get user code, sent by Instagram Login page
	urlParameters := r.URL.Query()
	code := urlParameters.Get("code")

	accessToken, err := authApi.ExchangeCodeForAccessToken(code)
	if err != nil {
		fmt.Fprintf(w, "Sorry an error occured %s", err)
		return
	}

	// Send cookie to client
	authCookie := &http.Cookie{
		Name:  "token",
		Value: accessToken,
	}
	http.SetCookie(w, authCookie)
}

/*
 Extract Access Token from request cookies
*/
func getTokenFromCookie(r *http.Request) string {
	cookie, _ := r.Cookie("token")
	return cookie.Value
}

func main() {
	// Create Instagram Login client
	authApi = insta.InstaLogin{ClientID: clientID, ClientSecret: clientSecret, RedirectURL: redirectURL}
	// Handle authentication
	http.HandleFunc("/login", RedirectToLogin)
	http.HandleFunc("/oauth-complete", GetAccessToken)

	http.HandleFunc("/search", SearchUser)

	http.ListenAndServe(":8080", nil)
}

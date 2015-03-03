package insta

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	//"golang.org/x/oauth2"
)

type TestUser struct {
	Name  string `json:"name"`
	Posts int    `json:"posts"`
}

func TestGetRequestParameters(t *testing.T) {
	options := map[string]string{
		"name":  "me",
		"count": "100",
	}

	// Dummy Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test GET parameters
		for k, v := range options {
			if sent := r.FormValue(k); sent != v {
				t.Errorf("Expected parameter for '%s': %s, got %s", k, v, sent)
			}
		}
	}))
	defer server.Close()

	c := NewClient(&http.Client{})
	c.getRequest(server.URL, options, nil)
}

func TestGetRequestInvalidURL(t *testing.T) {
	c := NewClient(&http.Client{})
	err := c.getRequest("kttp://abcd", make(map[string]string), nil)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if _, ok := err.(*url.Error); !ok {
		t.Errorf("Expected an URL error, got %#v", err)
	}
}

func TestDecodeBody(t *testing.T) {
	user := TestUser{
		Name:  "Adam",
		Posts: 100,
	}

	userString := `{"name":"Adam","posts":100}`

	var decodedUser TestUser
	err := decodeBody(bytes.NewBufferString(userString), &decodedUser)
	if err != nil {
		t.Errorf("Got error %s", err)
	}

	if decodedUser != user {
		t.Errorf("Expected %+v, got %+v", user, decodedUser)
	}
}

func TestAPIError(t *testing.T) {
	var meta ResponseMeta
	meta.Code = 400
	meta.ErrorType = "oauthexception"
	meta.ErrorMessage = "invalid access token"
	var r ApiResponse
	r.Meta = meta
	b, _ := json.Marshal(r)
	reader := bytes.NewBuffer(b)

	resp := http.Response{}
	resp.Body = ioutil.NopCloser(reader)

	apiError := newApiError(&resp)

	if !(meta.Code == apiError.Code &&
		meta.ErrorType == apiError.ErrorType &&
		meta.ErrorMessage == apiError.ErrorMessage) {
		t.Errorf("Expected %+v, got %+v", meta, apiError)
	}
}
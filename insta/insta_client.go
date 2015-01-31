package insta

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	// Get user's posts
	userFeedUrl = "https://api.instagram.com/v1/users/%s/media/recent/?access_token=%s"
	// Search for user based on username
	userSearchUrl = "https://api.instagram.com/v1/users/search?q=%s&count=1&access_token=%s"
)

// Instagram API client, it normally it requires an access
// token, but some parts of the API can be accessed by just
// using the client ID, please check the Instagram API doc
type InstaClient struct {
	ClientId    string
	AccessToken string
}

// Returns user ID associated with given username
func (i *InstaClient) GetUserId(username string) (string, error) {
	// Create URL for searching users with given username
	url := fmt.Sprintf(userSearchUrl, username, i.AccessToken)
	// Send GET request to find user
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// Decode JSON to get requested user
	decoder := json.NewDecoder(resp.Body)
	var result SearchResult
	err = decoder.Decode(&result)
	if err != nil {
		return "", err
	}
	// If no user was found, return
	if len(result.Users) == 0 {
		return "", errors.New("No user found with username " + username)
	}

	return result.Users[0].Id, nil
}

// Builds an URL for retrieving user's posts based on the options received
// it accepts the following arguments:
//
// - maxId: the retrieved posts will have an Id smaller than this
//
func (i *InstaClient) getPosts(username string, options map[string]string) ([]Post, error) {
	// Get User ID for given username
	userId, err := i.GetUserId(username)
	if err != nil {
		return nil, err
	}
	// Create generic URL to get user's posts
	url := fmt.Sprintf(userFeedUrl, userId, i.AccessToken)

	// Parse options
	if maxId, contains := options["maxId"]; contains {
		url += "&max_id=" + string(maxId)
	}

	return i.getPostsFromUrl(url)
}

// Returns the user's latest posts
func (i *InstaClient) GetRecentPosts(username string) ([]Post, error) {
	return i.getPosts(username, make(map[string]string))
}

// Returns a users's latest posts that have an ID smaller than maxId
func (i *InstaClient) GetPostsWithMaxId(username string, maxId string) ([]Post, error) {
	return i.getPosts(username, map[string]string{
		"maxId": maxId,
	})
}

// Retrieves and parses user's posts from given URL
func (i *InstaClient) getPostsFromUrl(url string) ([]Post, error) {
	// Get posts from given URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// Decode JSON to get posts
	decoder := json.NewDecoder(resp.Body)
	var result UserFeed
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	if len(result.Posts) == 0 {
		return nil, errors.New("This user doesn't have any public posts")
	}
	return result.Posts, nil
}

type SearchResult struct {
	Users []UserId `json:"data"`
}

// A single user's ID
type UserId struct {
	Id string `json:"id"`
}

// Instagram user's feed
type UserFeed struct {
	Posts []Post `json:"Data"`
}

// User post including image, likes, comments etc.
type Post struct {
	Id     string    `json:"id"`
	Images FeedImage `json:"images"`
}

// A single image with multiple resolutions
type FeedImage struct {
	Thumbnail          ImageUrl `json:"thumbnail"`
	LowResolution      ImageUrl `json:"low_resolution"`
	StandardResolution ImageUrl `json:"standard_resolution"`
}

// A single image url
type ImageUrl struct {
	Url string `json:"url"`
}

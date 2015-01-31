package insta

import (
	"errors"
)

// Instagram access token response
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// User search result
type SearchResult struct {
	Users []User `json:"data"`
}

// A single user
type User struct {
	Id string `json:"id"`
	Username string `json:"username"`
}

// Instagram user's feed
type UserFeed struct {
	Posts []Post `json:"Data"`
}

// Return Id of last post in feed
func (u *UserFeed) GetMinId() (string, error) {
	if len(u.Posts) == 0 {
		return "", errors.New("This feed contains no posts")
	}
	return u.Posts[len(u.Posts)-1].Id, nil
}

// User post including image, likes, comments etc.
type Post struct {
	Id       string        `json:"id"`
	Images   FeedImage     `json:"images"`
	Tags     []string      `json:"tags"`
	Caption  ImageCaption  `json:"caption"`
	Location ImageLocation `json:"location"`
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

// Caption for post
type ImageCaption struct {
	Text string `json:"text"`
}

type ImageLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

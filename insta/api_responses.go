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

// User struct with minimal information
// Passed along with every post when viewing a user's feed
type BaseUser struct {
	Id             string `json:"id"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}

// Passed along with every post when viewing currently
// authenticated user's liked posts or current user's feed
type UserWithFullName struct {
	BaseUser
	FullName string `json:"full_name"`
}

// Returned when searching for users based on query string
type UserWithStructuredName struct {
	BaseUser
	FirstName `json:"first_name"`
	LastName `json:"last_name"`
}

type UserWithFullDetails struct {
	UserWithFullName
	Bio string `json:"bio"`
	Website string `json:"website"`
	Counts UserStatistics `json:"counts"`
}

// Contains the number of posts, followers and follows of an user
type UserStatistics struct {
	Media int `json:"media"`
	Follows int `json:"follows"`
	FollowedBy int `json:"followed_by"`
}



// Instagram user's feed
type UserFeed struct {
	Posts []Post `json:"data"`
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
	Name      string  `json"location"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

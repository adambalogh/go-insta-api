package insta

import (
	"errors"
	"fmt"
)

const (
	// Get user's posts
	userFeedUrl = "users/%s/media/recent"
	// Search for user based on username
	userSearchUrl = "users/search?q=%s&count=1&access_token=%s"
)

// Returns user ID associated with given username
func (i *InstaClient) GetUserId(username string) (string, error) {
	var searchResult SearchResult
	i.get(userSearchUrl, map[string]string{
		"q": username,
	}, &searchResult)

	// If no user was found, return
	if len(searchResult.Users) == 0 {
		return "", errors.New("No user found with username " + username)
	}

	return searchResult.Users[0].Id, nil
}

// Builds an URL for retrieving user's posts based on the options received
// it accepts the following arguments:
//
// - max_id: the retrieved posts will have an Id smaller than this
//
func (i *InstaClient) getPosts(username string, options map[string]string) ([]Post, error) {
	// Get User ID for given username
	userId, err := i.GetUserId(username)
	if err != nil {
		return nil, err
	}
	// Create generic URL to get user's posts
	url := fmt.Sprintf(userFeedUrl, userId)

	var feed UserFeed
	err = i.get(url, options, &feed)
	if err != nil {
		return nil, err
	}

	return feed.Posts, nil
}

// Returns the user's latest posts
func (i *InstaClient) GetRecentPosts(username string) ([]Post, error) {
	return i.getPosts(username, make(map[string]string))
}

// Returns a users's latest posts that have an ID smaller than maxId
func (i *InstaClient) GetPostsWithMaxId(username string, maxId string) ([]Post, error) {
	return i.getPosts(username, map[string]string{
		"max_id": maxId,
	})
}

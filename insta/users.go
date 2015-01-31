package insta

import (
	"errors"
	"fmt"
)

// Search for user based on query string
func (i *InstaClient) SearchUser(queryString string, options map[string]string) (*SearchResult, error) {
	options["q"] = queryString

	var searchResult SearchResult
	err := i.get("/users/search", options, &searchResult)
	if err != nil {
		return nil, err
	}

	return &searchResult, nil
}

// Returns user ID associated with given username
func (i *InstaClient) GetUserId(username string) (string, error) {
	searchResult, err := i.SearchUser(username, map[string]string{"count": "1"})
	if err != nil {
		return "", err
	}

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
func (i *InstaClient) getPosts(userId string, options map[string]string) (*UserFeed, error) {
	var feed UserFeed
	err := i.get(fmt.Sprintf("/users/%s/media/recent", userId), options, &feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}

// Returns the user's latest posts
func (i *InstaClient) GetRecentPosts(userId string) (*UserFeed, error) {
	return i.getPosts(userId, make(map[string]string))
}

// Returns a users's latest posts that have an ID smaller than maxId
func (i *InstaClient) GetPostsWithMaxId(userId string, maxId string) (*UserFeed, error) {
	return i.getPosts(userId, map[string]string{
		"max_id": maxId,
	})
}

// Get currently logged in user's liked posts
func (i *InstaClient) GetLikedPosts(options map[string]string) (*UserFeed, error) {
	var feed UserFeed
	err := i.get("/users/self/media/liked", options, &feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}

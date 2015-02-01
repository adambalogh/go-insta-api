package insta

import (
	"errors"
	"fmt"
)

// Returns full user profile for user ID
func (i *InstaClient) GetUserProfile(userID string) (*UserWithFullDetails, error) {
	if len(userID) == 0 {
		return nil, errors.New("User ID cannot be empty")
	}

	var userProfileResult UserProfileResult
	err := i.requestURL(fmt.Sprintf("/users/%s", userID), map[string]string{}, &userProfileResult)
	if err != nil {
		return nil, err
	}
	return &userProfileResult.UserProfile, err
}

// Returns currently authenticated user's feed
func (i *InstaClient) GetSelfFeed() (*UserFeed, error) {
	var selfFeed UserFeed
	err := i.requestURL("/users/self/feed", map[string]string{}, &selfFeed)
	if err != nil {
		return nil, err
	}
	return &selfFeed, nil
}

// Searches for users based on query string
func (i *InstaClient) SearchUser(queryString string, options map[string]string) (*SearchResult, error) {
	if len(queryString) == 0 {
		return nil, errors.New("query string cannot be empty")
	}

	options["q"] = queryString // add query string to options map
	var searchResult SearchResult
	err := i.requestURL("/users/search", options, &searchResult)
	if err != nil {
		return nil, err
	}
	return &searchResult, nil
}

// Returns user ID associated with given username
func (i *InstaClient) GetUserID(username string) (string, error) {
	searchResult, err := i.SearchUser(username, map[string]string{"count": "1"})
	if err != nil {
		return "", err
	}
	// If no user was found, return
	if len(searchResult.Users) == 0 {
		return "", errors.New("No user found with username " + username)
	}
	return searchResult.Users[0].ID, nil
}

// Builds an URL for retrieving user's posts based on the options received
// it accepts the following arguments:
//
// - max_id: the retrieved posts will have an ID smaller than this
//
func (i *InstaClient) GetPosts(userID string, options map[string]string) (*UserFeed, error) {
	if len(userID) == 0 {
		return nil, errors.New("User ID cannot be empty")
	}

	var feed UserFeed
	err := i.requestURL(fmt.Sprintf("/users/%s/media/recent", userID), options, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

// Returns the user's latest posts
func (i *InstaClient) GetRecentPosts(userID string) (*UserFeed, error) {
	return i.GetPosts(userID, map[string]string{})
}

// Returns a users's latest posts that have an ID smaller than maxID
func (i *InstaClient) GetPostsWithMaxID(userID string, maxID string) (*UserFeed, error) {
	return i.GetPosts(userID, map[string]string{
		"max_id": maxID,
	})
}

// Gets the currently logged in user's liked posts
func (i *InstaClient) GetLikedPosts(options map[string]string) (*UserFeed, error) {
	var feed UserFeed
	err := i.requestURL("/users/self/media/liked", options, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

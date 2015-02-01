package insta

import (
	"errors"
	"fmt"
)

// Returns full user profile for user ID
func (i* InstaClient) GetUserProfile(userId string) (*UserWithFullDetails, error){
	var userProfileResult UserProfileResult
	err := i.requestUrl(fmt.Sprintf("/users/%s", userId), map[string]string{}, &userProfileResult)
	if err != nil {
		return nil, err
	}
	return &userProfileResult.UserProfile, err
}

// Returns currently authenticated user's feed
func (i *InstaClient) GetSelfFeed() (*UserFeed, error) {
	var selfFeed UserFeed
	err := i.requestUrl("/users/self/feed", map[string]string{}, &selfFeed)
	if err != nil {
		return nil, err
	}
	return &selfFeed, nil
}

// Searches for users based on query string
func (i *InstaClient) SearchUser(queryString string, options map[string]string) (*SearchResult, error) {
	options["q"] = queryString // add query string to options map
	var searchResult SearchResult
	err := i.requestUrl("/users/search", options, &searchResult)
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
func (i *InstaClient) GetPosts(userId string, options map[string]string) (*UserFeed, error) {
	var feed UserFeed
	err := i.requestUrl(fmt.Sprintf("/users/%s/media/recent", userId), options, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

// Returns the user's latest posts
func (i *InstaClient) GetRecentPosts(userId string) (*UserFeed, error) {
	return i.GetPosts(userId, map[string]string{})
}

// Returns a users's latest posts that have an ID smaller than maxId
func (i *InstaClient) GetPostsWithMaxId(userId string, maxId string) (*UserFeed, error) {
	return i.GetPosts(userId, map[string]string{
		"max_id": maxId,
	})
}

// Gets the currently logged in user's liked posts
func (i *InstaClient) GetLikedPosts(options map[string]string) (*UserFeed, error) {
	var feed UserFeed
	err := i.requestUrl("/users/self/media/liked", options, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

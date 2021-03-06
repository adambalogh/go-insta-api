package insta

import (
	"errors"
	"fmt"
	"sort"
)

// GetUserProfile returns the user profile with the given user ID.
func (i *InstaClient) GetUserProfile(userID string) (*UserWithFullDetails, error) {
	if len(userID) == 0 {
		return nil, errors.New("User ID cannot be empty")
	}

	userProfileResult := new(UserProfileResult)
	err := i.getRequest(base+fmt.Sprintf("/users/%s", userID), map[string]string{}, userProfileResult)
	if err != nil {
		return nil, err
	}
	return &userProfileResult.UserProfile, err
}

// GetSelfFeed returns the current user's feed.
func (i *InstaClient) GetSelfFeed() (*UserFeed, error) {
	selfFeed := new(UserFeed)
	err := i.getRequest(base+"/users/self/feed", map[string]string{}, selfFeed)
	if err != nil {
		return nil, err
	}
	return selfFeed, nil
}

// SearchUser searches for users based on the query string and returns the results.
func (i *InstaClient) SearchUser(queryString string, options map[string]string) (*SearchResult, error) {
	if len(queryString) == 0 {
		return nil, errors.New("query string cannot be empty")
	}

	options["q"] = queryString // add query string to options map
	searchResult := new(SearchResult)
	err := i.getRequest(base+"/users/search", options, searchResult)
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

// GetUserID returns the user ID associated with the requested username.
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

// GetPosts returns the requested user's posts.
// Any arguments for the API call should be in the options map.
func (i *InstaClient) GetPosts(userID string, options map[string]string) (*UserFeed, error) {
	if len(userID) == 0 {
		return nil, errors.New("User ID cannot be empty")
	}

	feed := new(UserFeed)
	err := i.getRequest(fmt.Sprintf(base+"/users/%s/media/recent", userID), options, feed)
	if err != nil {
		return nil, err
	}
	return feed, nil
}

// GetRecentPosts returns the requested user's latest posts
func (i *InstaClient) GetRecentPosts(userID string) (*UserFeed, error) {
	return i.GetPosts(userID, map[string]string{})
}

// GetPostsWithMaxID returns the requested users's latest posts
// that have an ID smaller than maxID
func (i *InstaClient) GetPostsWithMaxID(userID string, maxID string) (*UserFeed, error) {
	return i.GetPosts(userID, map[string]string{
		"max_id": maxID,
	})
}

// GetPostsFromUsers returns the merged feed of the requested users.
//
// The given options will be used for each individual user feed request,
// e.g. if options["count"] = 5, then it will return 5 posts from each of the
// requested users
func (i *InstaClient) GetPostsFromUsers(userIDs []string, options map[string]string) ([]Post, error) {
	var posts []Post
	var e error
	postsChannel := make(chan []Post)
	errorChannel := make(chan error)
	// Send requests for posts
	for _, userID := range userIDs {
		go func(userID string, options map[string]string, postsChannel chan []Post, errorChannel chan error) {
			feed, err := i.GetPosts(userID, options)
			if err != nil {
				errorChannel <- err
				return
			}
			postsChannel <- feed.Posts
		}(userID, options, postsChannel, errorChannel)
	}
	// Receive request results
	for i := 0; i < len(userIDs); i++ {
		select {
		case userPosts := <-postsChannel:
			posts = append(posts, userPosts...)
		case err := <-errorChannel:
			e = err
		}
	}
	// Sort posts by created time
	sort.Sort(ByCreatedTime(posts))

	return posts, e
}

// Sorting Interface for UserFeed.Post
type ByCreatedTime []Post

func (c ByCreatedTime) Len() int           { return len(c) }
func (c ByCreatedTime) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByCreatedTime) Less(i, j int) bool { return c[i].CreatedTime > c[j].CreatedTime }

// GetLikedPosts returns the currently logged in user's liked posts
func (i *InstaClient) GetLikedPosts(options map[string]string) (*UserFeed, error) {
	feed := new(UserFeed)
	err := i.getRequest(base+"/users/self/media/liked", options, feed)
	if err != nil {
		return nil, err
	}
	return feed, nil
}

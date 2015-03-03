package insta

import (
	"errors"
	"fmt"
)

// GetFollows returns the users that are followed by the requested user
func (i *InstaClient) GetFollows(userID string) ([]UserWithFullName, error) {
	if len(userID) == 0 {
		return nil, errors.New("User ID cannot be empty")
	}

	followsResponse := new(FollowsResult)
	err := i.getRequest(base+fmt.Sprintf("/users/%s/follows", userID), map[string]string{}, followsResponse)
	if err != nil {
		return nil, err
	}
	return followsResponse.Users, err
}

package insta

// Instagram access token response
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// Base struct for every Instagram API response
type ApiResponse struct {
	Meta ResponseMeta `json:"meta"`
}

// Metainfo for the response, can contain errors
type ResponseMeta struct {
	Code         int    `json:"code"`
	ErrorType    string `json:"error_type"`
	ErrorMessage string `json:"error_message"`
}

// Contains pagination info for sequential data
type PaginatedApiResponse struct {
	ApiResponse
	Pagination ResponsePagination
}

// Pagination info
type ResponsePagination struct {
	NextURL        string `json:"next_url"`
	NextMaxID      string `json:"next_max_id"`
	NextMaxLikedID string `json:"next_max_like_id"`
}

// User struct with minimal information
// Passed along with every post when viewing a user's feed
type BaseUser struct {
	ID             string `json:"id"`
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
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Returned when requesting a user's profile
type UserWithFullDetails struct {
	UserWithFullName
	Bio     string         `json:"bio"`
	Website string         `json:"website"`
	Counts  UserStatistics `json:"counts"`
}

// Contains the number of posts, followers and follows of an user
type UserStatistics struct {
	Media      int `json:"media"`
	Follows    int `json:"follows"`
	FollowedBy int `json:"followed_by"`
}

// FollowsResult contains a list of followed users
type FollowsResult struct {
	ApiResponse
	Users []UserWithFullName `json:"data"`
}

// User profile lookup result
type UserProfileResult struct {
	ApiResponse
	UserProfile UserWithFullDetails `json:"data"`
}

// User search result
type SearchResult struct {
	PaginatedApiResponse
	Users []UserWithStructuredName `json:"data"`
}

// Instagram user's feed
type UserFeed struct {
	PaginatedApiResponse
	Posts []Post `json:"data"`
}

// User post including image, likes, comments etc.
type Post struct {
	ID          string        `json:"id"`
	Images      PostImage     `json:"images"`
	Tags        []string      `json:"tags"`
	Caption     ImageCaption  `json:"caption"`
	Location    ImageLocation `json:"location"`
	CreatedTime string        `json:"created_time"`
}

// A single image with multiple resolutions
type PostImage struct {
	Thumbnail          ImageURL `json:"thumbnail"`
	LowResolution      ImageURL `json:"low_resolution"`
	StandardResolution ImageURL `json:"standard_resolution"`
}

// A single image url
type ImageURL struct {
	URL string `json:"url"`
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

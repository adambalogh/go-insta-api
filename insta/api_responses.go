package insta

// Instagram access token response
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type ApiMetaReporter interface {
	GetMeta() ResponseMeta
}

// Base struct for every Instagram API response
type ApiResponse struct {
	Meta ResponseMeta `json:"meta"`
}

func (a *ApiResponse) GetMeta() ResponseMeta{
	return a.Meta
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
	NextUrl       string `json:"next_url"`
	NextMaxId     string `json:"next_max_id"`
	NextMaxLikeId string `json:"next_max_like_id"`
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

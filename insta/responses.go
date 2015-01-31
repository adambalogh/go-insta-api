package insta

// Instagram access token response
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// User search result
type SearchResult struct {
	Users []UserId `json:"data"`
}

// A single user's ID
type UserId struct {
	Id string `json:"id"`
}

// Instagram user's feed
type UserFeed struct {
	Posts []Post `json:"Data"`
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

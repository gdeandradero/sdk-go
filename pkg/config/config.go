package config

// config is an mp instance.
var config *mp

// mp represents the config.
type mp struct {
	accessToken string
	productID   string
}

// New creates a new config.
func New(accessToken string) {
	config = &mp{
		accessToken: "Bearer " + accessToken,
		productID:   "123",
	}
}

// AccessToken returns the access token.
func AccessToken() string {
	return config.accessToken
}

// SetAccessToken sets the access token.
func SetAccessToken(at string) {
	config.accessToken = "Bearer " + at
}

// ProductID returns the product ID.
func ProductID() string {
	return config.productID
}

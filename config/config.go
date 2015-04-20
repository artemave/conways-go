package config

import "os"

// Config - app config
type Config struct{}

// GoogleClientID - returns google client id
func GoogleClientID() string {
	id := os.Getenv("GOOGLE_CLIENT_ID")
	if id == "" {
		panic("GOOGLE_CLIENT_ID is not set")
	}
	return id
}

// GoogleClientSecret - returns google client secret
func GoogleClientSecret() string {
	secret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if secret == "" {
		panic("GOOGLE_CLIENT_SECRET is not set")
	}
	return secret
}

// OauthRedirectURL - returns oauth redirect url
func OauthRedirectURL() string {
	url := os.Getenv("OAUTH_REDIRECT_URL")
	if url == "" {
		panic("OAUTH_REDIRECT_URL is not set")
	}
	return url
}

// RedisURL - returns redis address url
func RedisURL() string {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		panic("REDIS_URL is not set")
	}
	return url
}

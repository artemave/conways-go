package config

import "os"
import "net/url"

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

func redisURL() *url.URL {
	s := os.Getenv("REDIS_URL")
	if s == "" {
		panic("REDIS_URL is not set")
	}

	redisURL, err := url.Parse(s)

	if err != nil {
		panic("Bad password from redis url: " + s)
	}

	return redisURL
}

func RedisHost() string {
	return redisURL().Host
}

// RedisPassword - returns redis password if any
func RedisPassword() string {
	pass := ""
	if redisURL().User != nil {
		p, ok := redisURL().User.Password()
		if ok {
			pass = p
		}
	}
	return pass
}

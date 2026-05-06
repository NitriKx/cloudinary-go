// Package config defines the Cloudinary configuration.
package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/creasty/defaults"
	"github.com/gorilla/schema"
)

// Configuration is the main configuration struct.
type Configuration struct {
	Cloud     Cloud
	API       API
	URL       URL
	AuthToken AuthToken
}

var decoder = schema.NewDecoder()

// New returns a new Configuration instance from the environment variable
func New() (*Configuration, error) {
	return NewFromURL(os.Getenv("CLOUDINARY_URL"))
}

// NewFromURL returns a new Configuration instance from a cloudinary url.
func NewFromURL(cldURLStr string) (*Configuration, error) {
	if cldURLStr == "" {
		return nil, errors.New("must provide CLOUDINARY_URL")
	}

	cldURL, err := url.Parse(cldURLStr)
	if err != nil {
		return nil, err
	}

	pass, _ := cldURL.User.Password()
	params := cldURL.Query()
	conf, err := NewFromQueryParams(cldURL.Host, cldURL.User.Username(), pass, params)
	if err != nil {
		return nil, err
	}

	return conf, err
}

// NewFromAccountURL returns a new Configuration instance from a Cloudinary account URL.
//
// The expected format is:
//
//	account://<ACCOUNT_API_KEY>:<ACCOUNT_API_SECRET>@<ACCOUNT_ID>
//
// Where:
//   - ACCOUNT_API_KEY    is the account management API key (username part)
//   - ACCOUNT_API_SECRET is the account management API secret (password part)
//   - ACCOUNT_ID         is the account UUID (host part)
//
// This matches the CLOUDINARY_ACCOUNT_URL environment variable format from the Cloudinary CLI.
// Optional query parameters (e.g. upload_prefix, timeout) are forwarded to the configuration.
func NewFromAccountURL(accountURLStr string) (*Configuration, error) {
	if accountURLStr == "" {
		return nil, errors.New("must provide CLOUDINARY_ACCOUNT_URL")
	}

	u, err := url.Parse(accountURLStr)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "account" {
		return nil, fmt.Errorf("cloudinary account URL must have scheme 'account', got '%s'", u.Scheme)
	}

	pass, _ := u.User.Password()
	params := u.Query()
	// The host encodes the account ID; inject it as a query param so the
	// schema decoder populates Cloud.AccountID.
	params.Set("account_id", u.Host)
	return NewFromQueryParams("", u.User.Username(), pass, params)
}

// NewFromParams returns a new Configuration instance from the provided parameters.
func NewFromParams(cloud string, key string, secret string) (*Configuration, error) {
	return NewFromQueryParams(cloud, key, secret, map[string][]string{})
}

// NewFromOAuthToken returns a new Configuration instance from the provided cloud name and OAuth token.
func NewFromOAuthToken(cloud string, oAuthToken string) (*Configuration, error) {
	return NewFromQueryParams(cloud, "", "", map[string][]string{"oauth_token": {oAuthToken}})
}

// NewFromQueryParams returns a new Configuration instance from the provided url query parameters.
func NewFromQueryParams(cloud string, key string, secret string, params map[string][]string) (*Configuration, error) {
	cloudConf := Cloud{
		CloudName: cloud,
		APIKey:    key,
		APISecret: secret,
	}

	conf := &Configuration{
		Cloud:     cloudConf,
		API:       API{},
		URL:       URL{},
		AuthToken: AuthToken{},
	}

	if err := defaults.Set(conf); err != nil {
		return nil, err
	}

	// import configuration keys from parameters

	decoder.IgnoreUnknownKeys(true)

	err := decoder.Decode(&conf.Cloud, params)
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(&conf.API, params)
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(&conf.URL, params)
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(&conf.AuthToken, params)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/adetunjii/auth-svc/internal/port"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GOOGLE_API_URL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

var (
	ErrAuthenticationFailed = errors.New("user authentication faield")
)

type GoogleClient struct {
	logger port.AppLogger

	*oauth2.Config
}

type GoogleTokenDetails struct {
	accessToken  string
	expiry       string
	refreshToken string
}

func NewGoogleClient(client_id string, client_secret string, scopes []string, redirect_url string, logger port.AppLogger) *GoogleClient {

	oauth2Conf := &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		Scopes:       scopes,
		RedirectURL:  redirect_url,
		Endpoint:     google.Endpoint,
	}

	return &GoogleClient{Config: oauth2Conf, logger: logger}
}

func (g *GoogleClient) ExchangeCode(code string) (*GoogleTokenDetails, error) {

	// validate auth code
	if code == "" {
		return nil, errors.New("code not found to provide an access token")
	}

	// exchange auth code for token details
	token, err := g.Exchange(context.Background(), code)
	if err != nil {
		g.logger.Error("google code exchange failed with: ", err)
		return nil, errors.New("google code exchange failed")
	}

	gTokenDetails := &GoogleTokenDetails{
		accessToken:  token.AccessToken,
		expiry:       token.Expiry.String(),
		refreshToken: token.RefreshToken,
	}

	return gTokenDetails, nil
}

func (g *GoogleClient) FetchUserDetails(access_token string) (map[string]interface{}, error) {

	resp, err := http.Get(GOOGLE_API_URL + url.QueryEscape(access_token))
	if err != nil {
		g.logger.Error("user authentication failed with: ", err)
		return nil, ErrAuthenticationFailed
	}

	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		g.logger.Error("failed to read response body from google's server: ", err)
		return nil, ErrAuthenticationFailed
	}

	bytes := []byte(response)
	user := map[string]interface{}{}

	if err := json.Unmarshal(bytes, &user); err != nil {
		g.logger.Error("failed to unmarshal google user struct", err)
		return nil, errors.New("failed to unmarshal user object")
	}

	return user, nil
}

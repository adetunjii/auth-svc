package oauth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/adetunjii/auth-svc/internal/port"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type FacebookClient struct {
	logger port.AppLogger

	*oauth2.Config
}

var FACEBOOK_API_URL = "https://graph.facebook.com/me?fields=id,name,email&access_token="

func NewFaceBookClient(client_id string, client_secret string, scopes []string, redirect_url string, logger port.AppLogger) *FacebookClient {
	oauth2Conf := &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		Scopes:       scopes,
		RedirectURL:  redirect_url,
		Endpoint:     facebook.Endpoint,
	}

	return &FacebookClient{Config: oauth2Conf, logger: logger}
}

func (f *FacebookClient) FetchUserDetails(access_token string) (map[string]interface{}, error) {
	resp, err := http.Get(FACEBOOK_API_URL + url.QueryEscape(access_token))
	if err != nil {
		f.logger.Error("user authentication failed with: ", err)
		return nil, ErrAuthenticationFailed
	}

	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		f.logger.Error("failed to read response body from facebook: ", err)
		return nil, ErrAuthenticationFailed
	}

	bytes := []byte(response)
	user := map[string]interface{}{}

	if err := json.Unmarshal(bytes, &user); err != nil {
		f.logger.Error("failed to unmarshal facebook user struct with err: ", err)
		return nil, errors.New("failed to unmarshal user object")
	}

	return user, nil
}

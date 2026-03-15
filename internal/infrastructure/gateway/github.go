package gateway

import (
	"auth-api/internal/domain"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type GithubApi interface {
	GetAppToken(code string) (string, error)
	GetUserInfo(code string) (domain.User, error)
}

type githubApi struct {
	clientId     string
	clientSecret string
	client       *http.Client
}

var _ GithubApi = (*githubApi)(nil)

func NewGithubApi(clientId, clientSecret string, client *http.Client) GithubApi {
	return &githubApi{clientId: clientId, clientSecret: clientSecret, client: client}
}

func (g githubApi) GetAppToken(code string) (string, error) {
	params := map[string]string{
		"client_id":     g.clientId,
		"client_secret": g.clientSecret,
		"code":          code,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var appTokenResp struct {
		AccessToken string  `json:"access_token"`
		Error       *string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&appTokenResp); err != nil {
		return "", err
	}

	if appTokenResp.Error != nil {
		return "", errors.New(*appTokenResp.Error)
	}

	return appTokenResp.AccessToken, nil
}

func (g githubApi) GetUserInfo(code string) (domain.User, error) {
	token, err := g.GetAppToken(code)

	if err != nil {
		return domain.User{}, err
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return domain.User{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := g.client.Do(req)
	if err != nil {
		return domain.User{}, err
	}
	defer resp.Body.Close()

	var githubUser struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		ID:    strconv.Itoa(githubUser.ID),
		Name:  githubUser.Name,
		Email: githubUser.Email,
	}

	return user, nil
}

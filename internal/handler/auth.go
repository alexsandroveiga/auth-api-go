package handler

import (
	"auth-api-go/internal/httpx"
	"auth-api-go/internal/infrastructure/gateway"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthHandler interface {
	LoginWithGithub(w http.ResponseWriter, r *http.Request) error
}

type authHandler struct {
	githubApi gateway.GithubApi
	token     gateway.JwtTokenHandler
}

var _ AuthHandler = (*authHandler)(nil)

func NewAuthHandler(githubApi gateway.GithubApi, token gateway.JwtTokenHandler) AuthHandler {
	return &authHandler{githubApi, token}
}

type LoginWithGithubRequest struct {
	Code string `json:"code"`
}

func (h *authHandler) LoginWithGithub(w http.ResponseWriter, r *http.Request) error {
	var request LoginWithGithubRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return err
	}

	code := request.Code
	user, err := h.githubApi.GetUserInfo(code)

	fmt.Println(user)

	if err != nil {
		return httpx.Unauthorized(err.Error())
	}

	token, err := h.token.Generate(user.ID, 3600000)

	if err != nil {
		return httpx.Unauthorized(err.Error())
	}

	auth_token := map[string]string{
		"token": token,
	}
	return httpx.JSON(w, 200, auth_token)
}

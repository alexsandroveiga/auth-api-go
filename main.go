package main

import (
	"auth-api-go/internal/handler"
	"auth-api-go/internal/httpx"
	"auth-api-go/internal/infrastructure/gateway"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	GITHUB_CLIENT_ID     = "GITHUB_CLIENT_ID"
	GITHUB_CLIENT_SECRET = "GITHUB_CLIENT_SECRET"
	JWT_SECRET           = "JWT_SECRET"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()

	clientId := os.Getenv(GITHUB_CLIENT_ID)
	clientSecret := os.Getenv(GITHUB_CLIENT_SECRET)
	client := &http.Client{}
	githubApi := gateway.NewGithubApi(clientId, clientSecret, client)

	jwtSecret := os.Getenv(JWT_SECRET)
	token := gateway.NewJwtTokenHandler(jwtSecret)

	userHandler := handler.NewUserHandler()
	authHandler := handler.NewAuthHandler(githubApi, token)

	mux.HandleFunc("POST /users", httpx.Adapt(userHandler.Create))
	mux.HandleFunc("GET /users", httpx.Adapt(userHandler.List))

	mux.HandleFunc("POST /login/github", httpx.Adapt(authHandler.LoginWithGithub))

	fmt.Println("Server running on :3030")

	http.ListenAndServe(os.Getenv("PORT"), mux)
}

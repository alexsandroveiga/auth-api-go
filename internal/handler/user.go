package handler

import (
	"auth-api/internal/domain"
	"auth-api/internal/httpx"
	"encoding/json"
	"net/http"
)

func Find[T any](items []T, fn func(T) bool) *T {
	for _, item := range items {
		if fn(item) {
			return &item
		}
	}
	return nil
}

type UserHandler interface {
	Create(w http.ResponseWriter, r *http.Request) error
	List(w http.ResponseWriter, r *http.Request) error
}

type userHandler struct {
	users []domain.User
}

var _ UserHandler = (*userHandler)(nil)

func NewUserHandler() UserHandler {
	return &userHandler{}
}

func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var user domain.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return err
	}

	userExists := Find(h.users, func(u domain.User) bool {
		return u.Email == user.Email
	})

	if userExists != nil {
		return httpx.UnprocessableEntity("user already exists")
	}

	h.users = append(h.users, user)

	return httpx.JSON(w, http.StatusCreated, user)
}

func (h *userHandler) List(w http.ResponseWriter, r *http.Request) error {
	users := h.users
	if users == nil {
		users = []domain.User{}
	}

	return httpx.JSON(w, http.StatusOK, users)
}

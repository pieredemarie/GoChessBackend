package auth

import "backendChess/pkg/storage"

type Handler struct {
	Storage *storage.PostgresStorage
}

type ErrorResponce struct {
	Message string `json:"message"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type LoginResponce struct {
	Token string `json:"token"`
}

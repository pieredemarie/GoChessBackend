package storage

import chesslogic "backendChess/pkg/chessLogic"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserInfo struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Elo        int    `json:"rating"`
	GamesCount int    `json:"count"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthStorage interface {
	Login(email, password string) (string, error)
	Register(data User) error
	GetUserInfo(userID int) (*UserInfo, error)
}

type GameStorage interface {
	SaveGame(game *chesslogic.Game) error 
	SaveMode(gameID string,move chesslogic.Move) error 
	UpdateGameResult(gameID string, winner string, reason string) error
}
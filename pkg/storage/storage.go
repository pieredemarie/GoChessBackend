package storage

import (
	chesslogic "backendChess/pkg/chessLogic"
	"time"
)

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

type DBGame struct {
	ID string 
	WhiteID int 
	BlackID int 
	Status string // "active/finished"
	CreatedAt time.Time 
	FinishedAt *time.Time
	Result string
}

type DBMove struct {
	ID string 
	GameID string 
	MoveNum int 
	From string 
	To string 
	Piece string 
	CreatedAt time.Time
}

type AuthStorage interface {
	Login(email, password string) (string, error)
	Register(data User) error
	GetUserInfo(userID int) (*UserInfo, error)
}

type GameStorage interface {
	SaveGame(g *DBGame) error 
	SaveMode(gameID string,move chesslogic.Move) error 
	UpdateGameResult(gameID string, winner string, reason string) error
}
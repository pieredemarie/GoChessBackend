package storage

import (
	"backendChess/pkg/jwtutils"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db: db,
	}, nil	
}

func (p *PostgresStorage) Register(data User) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(data.Password),bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_,err = p.db.Exec("INSERT into users (username, email, password_hash) VALUES ($1,$2,$3)", data.Username,data.Email,hashedPass)
	return err
}

func (p *PostgresStorage) Login(email, password string) (string, error) {
	var hashedPwd string 
	var userID int
	err := p.db.QueryRow("SELECT id, password_hash FROM users WHERE email = $1", email).Scan(&userID,&hashedPwd)
	if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return "", fmt.Errorf("user not found")
        }
        return "", fmt.Errorf("database error: %w", err)
    }

	isPasswordCorrect := bcrypt.CompareHashAndPassword([]byte(hashedPwd),[]byte(password))
	if isPasswordCorrect != nil {
		return "", fmt.Errorf("invalid password")
	}

	secret := os.Getenv("JWT_SECRET")
	token, err := jwtutils.GenerateToken(userID, secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (p *PostgresStorage) GetUserInfo(userId int) (*UserInfo, error) {
	var info UserInfo 

	err := p.db.QueryRow("SELECT id, username,rating,games_count FROM users WHERE id = $1",userId).Scan(&info.ID,&info.Username,&info.Elo,&info.GamesCount)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
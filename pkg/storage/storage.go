package storage

type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserInfo struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Elo int `json:"rating"`
	GamesCount int `json:"count"`
}

type Login struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type AuthStorage interface {
	Login(email, password string) (string, error)
	Register(data User) error
	GetUserInfo(userID int) (*UserInfo, error)
}
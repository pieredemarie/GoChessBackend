package storage

type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type AuthStorage interface {
	Login(email, password string) (string, error)
	Register(data User) error
}
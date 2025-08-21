package storage

import "database/sql"

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
	hashedPass := "here will be bcrypt yea"
	_,err := p.db.Exec("INSERT into users (username, email, password) VALUES ($1,$2,$3)", data.Username,data.Email,hashedPass)
	return err
}
package storage

import (
	"database/sql"
	"eWallet/config"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type DBStorage struct {
	conn *sql.DB
}

func NewDBStorage(config config.Config) (*DBStorage, error) {
	db, err := sql.Open("postgres", config.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db %w", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate driver, %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"ewallet", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("failed to do migrate %w", err)
	}
	s := &DBStorage{
		conn: db,
	}

	return s, s.CheckConnection()
}

func (s *DBStorage) CheckConnection() error {
	if err := s.conn.Ping(); err != nil {
		return fmt.Errorf("failed to connect to db %w", err)
	}
	return nil
}

func (s *DBStorage) Close() error {
	return s.conn.Close()
}

func (s *DBStorage) SaveWallet(id string, balance float64) error {
	insertQuery := `INSERT INTO wallets (idOfWallet, balance) 
		VALUES ($1, $2)
		ON CONFLICT (idOfWallet) 
		DO UPDATE SET balance = $2`
	_, err := s.conn.Exec(insertQuery, id, balance)

	if err != nil {
		return fmt.Errorf("failed to save id and balance in db: %v", err)
	}
	return nil
}

func (s *DBStorage) TakeWallet(id string) (string, float64, error) {
	row := s.conn.QueryRow("SELECT idOfWallet, balance FROM wallets WHERE idOfWallet = $1", id)
	var idOfWallet string
	var balance float64
	err := row.Scan(&idOfWallet, &balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", 0.0, fmt.Errorf("no idOfWallet from person")
		} else {
			return "", 0.0, fmt.Errorf("error: %w", err)
		}
	}
	return idOfWallet, balance, nil
}

func (s *DBStorage) UpdateWallet(id string, balance float64) error {
	_, err := s.conn.Exec("UPDATE wallets SET balance = balance + $1 WHERE idofwallet = $2", balance, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no idOfWallet to person")
		} else {
			return fmt.Errorf("didn't update balance: %w", err)
		}
	}
	return nil
}

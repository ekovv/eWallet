package storage

import (
	"context"
	"database/sql"
	"eWallet/config"
	"eWallet/internal/shema"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type DBStorage struct {
	conn *sql.DB
}

func NewDBStorage(config config.Config) (*DBStorage, error) {
	if config.TypeDB == "postgresql" {
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
	} else {
		db, err := sql.Open("mongodb", config.DB)
		mongoDB, err := mongodb.WithInstance(db, &mongodb.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to create migrate driver, %w", err)
		}

		// Создание нового экземпляра миграции
		m, err := migrate.NewWithDatabaseInstance(
			"file://migrations",
			"mongodb", mongoDB)
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

func (s *DBStorage) SaveWallet(ctx context.Context, id string, balance float64) error {
	insertQuery := `INSERT INTO wallets (idOfWallet, balance) VALUES ($1, $2) ON CONFLICT (idOfWallet) DO UPDATE SET balance = $2`
	_, err := s.conn.ExecContext(ctx, insertQuery, id, balance)

	if err != nil {
		return fmt.Errorf("failed to save id and balance in db: %v", err)
	}
	return nil
}

func (s *DBStorage) TakeWallet(ctx context.Context, id string) (string, float64, error) {
	row := s.conn.QueryRowContext(ctx, "SELECT idOfWallet, balance FROM wallets WHERE idOfWallet = $1", id)
	var idOfWallet string
	var balance float64
	err := row.Scan(&idOfWallet, &balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", 0.0, fmt.Errorf("failed to get from person: %w", err)
		} else {
			return "", 0.0, fmt.Errorf("error: %w", err)
		}
	}
	return idOfWallet, balance, nil
}

func (s *DBStorage) UpdateWallet(ctx context.Context, id string, balance float64) error {
	res, err := s.conn.ExecContext(ctx, "UPDATE wallets SET balance = balance + $1 WHERE idofwallet = $2", balance, id)
	if err != nil {
		return fmt.Errorf("didn't update balance: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve the number of rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("failed to get to person: %w", err)
	}

	return nil
}

func (s *DBStorage) SaveInfo(ctx context.Context, from string, to string, amount float64, time string) error {
	insertQuery := `INSERT INTO history (time, from_wallet, to_wallet, amount) VALUES ($1, $2, $3, $4)`
	_, err := s.conn.ExecContext(ctx, insertQuery, time, from, to, amount)
	if err != nil {
		return fmt.Errorf("didn't save info: %w", err)
	}
	return nil
}

func (s *DBStorage) GetInfo(ctx context.Context, id string) ([]shema.HistoryTransfers, error) {
	rows, err := s.conn.QueryContext(ctx, "SELECT time, from_wallet, to_wallet, amount FROM history WHERE from_wallet = $1 OR to_wallet = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get from person: %w", err)
		} else {
			return nil, fmt.Errorf("didn't get info: %w", err)
		}
	}
	defer rows.Close()

	var history []shema.HistoryTransfers
	for rows.Next() {
		var t shema.HistoryTransfers
		err := rows.Scan(&t.Time, &t.From, &t.To, &t.Amount)
		if err != nil {
			return nil, err
		}
		history = append(history, t)
	}
	if history == nil {
		return nil, fmt.Errorf("failed to get from person: %w", err)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return history, nil
}

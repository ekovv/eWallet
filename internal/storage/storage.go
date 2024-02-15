package storage

import (
	"context"
	"eWallet/config"
	"eWallet/internal/shema"
	"eWallet/internal/storage/postgresql"
	"fmt"
)

type Storage interface {
	CheckConnection() error
	Close() error
	SaveWallet(ctx context.Context, id string, balance float64) error
	TakeWallet(ctx context.Context, id string) (string, float64, error)
	UpdateWallet(ctx context.Context, id string, balance float64) error
	SaveInfo(ctx context.Context, from string, to string, amount float64, time string) error
	GetInfo(ctx context.Context, id string) ([]shema.HistoryTransfers, error)
}

func New(cfg config.Config) (Storage, error) {
	if cfg.TypeDB == "postgresql" {
		d, err := postgresql.NewPostgresStorage(cfg)
		if err != nil {
			return nil, fmt.Errorf("error creating db storage: %v", err)
		}
		return d, nil
	} else {
		d, err := postgresql.NewPostgresStorage(cfg)
		if err != nil {
			return nil, fmt.Errorf("error creating db storage: %v", err)
		}
		return d, nil
	}
}

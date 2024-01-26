package service

import (
	"context"
	"eWallet/config"
	"eWallet/internal/constants"
	"eWallet/internal/domains"
	"eWallet/internal/shema"
	"fmt"
	"github.com/speps/go-hashids/v2"
	"go.uber.org/zap"
	"math/rand"
	"strings"
	"time"
)

type Service struct {
	storage domains.Storage
	config  config.Config
	logger  *zap.Logger
}

func NewService(storage domains.Storage, config config.Config) *Service {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil
	}
	return &Service{storage: storage, logger: logger, config: config}
}

func (s *Service) GenerateWallet(ctx context.Context) (string, float64, error) {
	const log = "Service.GenerateWallet"
	id, err := s.GenerateID()
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s : %v", log, err))
		return "", 0.0, err
	}
	balance := 100.0
	err = s.storage.SaveWallet(id, balance, ctx)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s : %v", log, err))
		return "", 0.0, err
	}
	return id, balance, nil
}

func (s *Service) GenerateID() (string, error) {
	const log = "Service.GenerateID"
	hd := hashids.NewData()
	hd.Salt = s.config.Salt
	hd.MinLength = 30
	h, err := hashids.NewWithData(hd)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s : %v", log, err))
		return "", err
	}

	id := []int{rand.Intn(1000000000)}

	hash, err := h.Encode(id)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s : %v", log, err))
		return "", err
	}

	return hash, nil
}

func (s *Service) Transaction(from string, to string, amount float64, ctx context.Context) error {
	const log = "Service.Transaction"
	idOfWallet, balance, err := s.storage.TakeWallet(from, ctx)
	if err != nil {
		if strings.Contains(err.Error(), "failed to get from person") {
			s.logger.Info(fmt.Sprintf("%s : %v", log, err))
			return constants.ErrNotFromPerson
		} else {
			s.logger.Info(fmt.Sprintf("%s : %v", log, err))
			return fmt.Errorf("didn't take from db: %w", err)
		}
	}

	if balance-amount < 0 {
		return constants.ErrBadAmount
	}

	balance = balance - amount
	err = s.storage.SaveWallet(idOfWallet, balance, ctx)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s : %v", log, err))
		return fmt.Errorf("didn't save new balance for from person: %w", err)
	}

	err = s.storage.UpdateWallet(to, amount, ctx)
	if err != nil {
		if strings.Contains(err.Error(), "failed to get to person") {
			s.logger.Info(fmt.Sprintf("%s : %v", log, err))
			return constants.ErrNotToPerson
		} else {
			s.logger.Info(fmt.Sprintf("%s : %v", log, err))
			return fmt.Errorf("didn't update to db: %w", err)
		}
	}
	t := time.Now()
	err = s.storage.SaveInfo(from, to, amount, t.Format(time.RFC3339), ctx)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s : %v", log, err))
		return fmt.Errorf("didn't save history in db: %w", err)
	}

	return nil
}

func (s *Service) GetHistory(id string, ctx context.Context) ([]shema.HistoryTransfers, error) {
	const log = "Service.GetHistory"
	history, err := s.storage.GetInfo(id, ctx)
	if err != nil {
		if strings.Contains(err.Error(), "failed to get from person") {
			s.logger.Info(fmt.Sprintf("%s : %v", log, err))
			return nil, constants.ErrNotFromPerson
		} else {
			s.logger.Info(fmt.Sprintf("%s : %v", log, err))
			return nil, fmt.Errorf("didn't take history from db: %w", err)
		}
	}
	return history, nil
}

func (s *Service) GetStatus(id string, ctx context.Context) (string, float64, error) {
	const log = "Service.GetStatus"
	idOfWallet, balance, err := s.storage.TakeWallet(id, ctx)
	if err != nil {
		if strings.Contains(err.Error(), "failed to get from person") {
			s.logger.Info(fmt.Sprintf("%s : %v", log, err))
			return "", 0.0, constants.ErrNotFromPerson
		} else {
			s.logger.Info(fmt.Sprintf("%s : %v", log, err))
			return "", 0.0, fmt.Errorf("didn't take from db: %w", err)
		}
	}
	return idOfWallet, balance, nil
}

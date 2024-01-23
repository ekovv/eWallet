package service

import (
	"eWallet/internal/constants"
	"eWallet/internal/domains"
	"eWallet/internal/shema"
	"errors"
	"fmt"
	"github.com/speps/go-hashids/v2"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

type Service struct {
	storage domains.Storage
	logger  *zap.Logger
}

func NewService(storage domains.Storage) *Service {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil
	}
	return &Service{storage: storage, logger: logger}
}

func (s *Service) GenerateWallet() (string, float64, error) {
	id, err := s.GenerateID()
	if err != nil {
		s.logger.Info("can't create id")
		return "", 0.0, err
	}
	balance := 100.0
	err = s.storage.SaveWallet(id, balance)
	if err != nil {
		s.logger.Info("didn't save in db")
		return "", 0.0, err
	}
	return id, balance, nil
}

func (s *Service) GenerateID() (string, error) {
	hd := hashids.NewData()
	hd.Salt = "my salt"
	hd.MinLength = 30
	h, err := hashids.NewWithData(hd)
	if err != nil {
		s.logger.Info("can't create HashID")
		return "", err
	}

	id := []int{rand.Intn(1000000000)}

	hash, err := h.Encode(id)
	if err != nil {
		s.logger.Info("can't encode hash")
		return "", err
	}

	return hash, nil
}

func (s *Service) Transaction(from string, to string, amount float64) error {
	idOfWallet, balance, err := s.storage.TakeWallet(from)
	if err != nil {
		if errors.Is(err, constants.ErrNotFromPerson) {
			s.logger.Info("no idOfWallet from person")
			return constants.ErrNotFromPerson
		} else {
			s.logger.Info("didn't take from db")
			return fmt.Errorf("didn't take from db: %w", err)
		}
	}

	t := time.Now()
	err = s.storage.SaveInfo(from, to, amount, t.Format(time.RFC3339))
	if err != nil {
		s.logger.Info("didn't save history in db")
		return fmt.Errorf("didn't save history in db: %w", err)
	}

	if balance-amount < 0 {
		return constants.ErrBadAmount
	}

	balance = balance - amount
	err = s.storage.SaveWallet(idOfWallet, balance)
	if err != nil {
		s.logger.Info("didn't save new balance for from person")
		return fmt.Errorf("didn't save new balance for from person: %w", err)
	}

	err = s.storage.UpdateWallet(to, amount)
	if err != nil {
		if errors.Is(err, constants.ErrNotToPerson) {
			s.logger.Info("no idOfWallet to person")
			return constants.ErrNotToPerson
		} else {
			s.logger.Info("didn't update to db")
			return fmt.Errorf("didn't update to db: %w", err)
		}
	}

	return nil
}

func (s *Service) GetHistory(id string) ([]shema.HistoryTransfers, error) {
	history, err := s.storage.GetInfo(id)
	if err != nil {
		if errors.Is(err, constants.ErrNotFromPerson) {
			s.logger.Info("no idOfWallet from person")
			return nil, constants.ErrNotFromPerson
		} else {
			s.logger.Info("didn't take history from db")
			return nil, fmt.Errorf("didn't take history from db: %w", err)
		}
	}
	return history, nil
}

func (s *Service) GetStatus(id string) (string, float64, error) {
	idOfWallet, balance, err := s.storage.TakeWallet(id)
	if err != nil {
		if errors.Is(err, constants.ErrNotFromPerson) {
			s.logger.Info("no idOfWallet from person")
			return "", 0.0, constants.ErrNotFromPerson
		} else {
			s.logger.Info("didn't take from db")
			return "", 0.0, fmt.Errorf("didn't take from db: %w", err)
		}
	}
	return idOfWallet, balance, nil
}

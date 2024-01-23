package service

import (
	"eWallet/internal/constants"
	"eWallet/internal/domains"
	"errors"
	"fmt"
	"github.com/speps/go-hashids/v2"
	"go.uber.org/zap"
	"math/rand"
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
		if errors.Is(err, fmt.Errorf("no idOfWallet from person")) {
			s.logger.Info("no idOfWallet from person")
			return constants.ErrNotFromPerson
		} else {
			s.logger.Info("didn't take from db")
			return fmt.Errorf("didn't take from db: %w", err)
		}
	}
	if balance-amount < 0 {
		return fmt.Errorf("bad amount")
	}
	balance = balance - amount

	err = s.storage.SaveWallet(idOfWallet, balance)
	if err != nil {
		s.logger.Info("didn't save new balance for from person")
		return fmt.Errorf("didn't save new balance for from person: %w", err)
	}
	err = s.storage.UpdateWallet(to, amount)
	if err != nil {
		if errors.Is(err, fmt.Errorf("no idOfWallet to person")) {
			s.logger.Info("no idOfWallet to person")
			return constants.ErrNotToPerson
		} else {
			s.logger.Info("didn't update to db")
			return fmt.Errorf("didn't update to db: %w", err)
		}
	}
	return nil
}

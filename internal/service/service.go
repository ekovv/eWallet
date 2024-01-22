package service

import (
	"eWallet/internal/domains"
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

func (s *Service) GenerateWallet() (string, error) {
	id, err := s.GenerateID()
	if err != nil {
		s.logger.Info("can't create id")
		return "", err
	}
	balance := 100.0
}

func (s *Service) GenerateID() (string, error) {
	hd := hashids.NewData()
	hd.Salt = "my salt"
	hd.MinLength = 8
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

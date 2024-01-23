package service

import (
	"eWallet/internal/constants"
	"eWallet/internal/domains/mocks"
	"errors"
	"go.uber.org/zap"
	"testing"
)

type storageMock func(c *mocks.Storage)

func TestService_Transaction(t *testing.T) {
	type args struct {
		from   string
		to     string
		amount float64
	}
	tests := []struct {
		name        string
		args        args
		storageMock storageMock
		wantErr     error
	}{
		{
			name: "OK1",
			args: args{
				from:   "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn",
				to:     "0MEe1ArlKnX5ea0ByX85PD83QBwpJa",
				amount: 20.0,
			},

			storageMock: func(c *mocks.Storage) {
				c.Mock.On("TakeWallet", "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn").Return("Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn", 100.0, nil).Times(1)
				c.Mock.On("SaveWallet", "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn", 80.0).Return(nil).Times(1)
				c.Mock.On("UpdateWallet", "0MEe1ArlKnX5ea0ByX85PD83QBwpJa", 20.0).Return(nil).Times(1)
			},
			wantErr: nil,
		},
		{
			name: "BAD",
			args: args{
				from:   "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn",
				to:     "0MEe1ArlKnX5ea0ByX85PD83ahsyu2",
				amount: 20.0,
			},

			storageMock: func(c *mocks.Storage) {
				c.Mock.On("TakeWallet", "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn").Return("Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn", 100.0, nil).Times(1)
				c.Mock.On("SaveWallet", "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn", 80.0).Return(nil).Times(1)
				c.Mock.On("UpdateWallet", "0MEe1ArlKnX5ea0ByX85PD83ahsyu2", 20.0).Return(constants.ErrNotToPerson).Times(1)
			},
			wantErr: constants.ErrNotToPerson,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := mocks.NewStorage(t)
			tt.storageMock(storage)
			logger, err := zap.NewProduction()

			service := Service{
				storage: storage,
				logger:  logger,
			}
			err = service.Transaction(tt.args.from, tt.args.to, tt.args.amount)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

package handler

import (
	"bytes"
	"eWallet/config"
	"eWallet/internal/constants"
	"eWallet/internal/domains/mocks"
	"eWallet/internal/shema"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type serviceMock func(c *mocks.Service)

func TestHandler_CreateWallet(t *testing.T) {
	tests := []struct {
		name        string
		serviceMock serviceMock
		wantCode    int
		want        shema.Wallet
	}{
		{
			name: "OK#1",
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GenerateWallet", mock.Anything).Return("BXbzKE4P62ui6evzzeB5wYLoO1r0Al", 100.0, nil).Times(1)
			},
			wantCode: http.StatusOK,
			want:     shema.Wallet{ID: "BXbzKE4P62ui6evzzeB5wYLoO1r0Al", Balance: 100},
		},
		{
			name: "OK#2",
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GenerateWallet", mock.Anything).Return("BXbzKE4P62ui6qwereB5wYLoO1r0Al", 100.0, nil).Times(1)
			},
			wantCode: http.StatusOK,
			want:     shema.Wallet{ID: "BXbzKE4P62ui6qwereB5wYLoO1r0Al", Balance: 100},
		},
		{
			name: "BAD",
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GenerateWallet", mock.Anything).Return("", 0.0, errors.New("can't create id")).Times(1)
			},
			wantCode: http.StatusBadRequest,
			want:     shema.Wallet{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			service := mocks.NewService(t)
			h := NewHandler(service, config.Config{})
			tt.serviceMock(service)

			path := "/api/v1/wallet"
			g.POST(path, h.CreateWallet)
			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, path, nil)

			g.ServeHTTP(w, request)

			if w.Code != tt.wantCode {
				t.Errorf("got %d, want %d", w.Code, tt.wantCode)
			}

			wantResponse, err := json.Marshal(tt.want)
			if err != nil {
				fmt.Errorf("failed json")
				return
			}

			response, err := json.Marshal(w.Body)
			if err != nil {
				fmt.Errorf("failed json")
				return
			}

			if bytes.Equal(wantResponse, response) {
				t.Errorf("got %s, want %v", w.Body, tt.want)
			}
		})
	}
}

func TestHandler_Transactions(t *testing.T) {
	type jso struct {
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}
	tests := []struct {
		name        string
		from        string
		json        jso
		serviceMock serviceMock
		wantCode    int
	}{
		{
			name: "OK1",
			from: "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn",
			json: jso{
				To:     "0MEe1ArlKnX5ea0ByX85PD83QBwpJa",
				Amount: 20.0,
			},
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("Transaction", mock.Anything, "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn", "0MEe1ArlKnX5ea0ByX85PD83QBwpJa", 20.0).Return(nil).Times(1)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "BAD",
			from: "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn",
			json: jso{
				To:     "jddjkjdksjdkjskd",
				Amount: 20.0,
			},
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("Transaction", mock.Anything, "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn", "jddjkjdksjdkjskd", 20.0).Return(fmt.Errorf("no idOfWallet to person")).Times(1)
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()

			service := mocks.NewService(t)
			h := NewHandler(service, config.Config{})

			path := "/api/v1/wallet/:walletId/send"
			g.POST(path, h.Transactions)
			tt.serviceMock(service)
			b, err := json.Marshal(tt.json)
			if err != nil {
				return
			}
			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/wallet/%s/send", tt.from), strings.NewReader(string(b)))

			g.ServeHTTP(w, request)

			if w.Code != tt.wantCode {
				t.Errorf("got %d, want %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestHandler_History(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		serviceMock serviceMock
		wantCode    int
		want        []shema.HistoryTransfers
	}{
		{
			name: "OK1",
			id:   "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn",
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GetHistory", mock.Anything, "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn").Return(
					[]shema.HistoryTransfers{
						{
							Time:   "2024-01-23 13:50:19.000000 +00:00",
							From:   "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn",
							To:     "aqaWKX3YeQ6jRWxJ32Zbxwq17kAdEn",
							Amount: 20.0,
						},
						{
							Time:   "2024-01-23 13:50:58.000000 +00:00",
							From:   "ay7uwX3YeQ6jRWxJ32Zbxwq17kAdEn",
							To:     "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn",
							Amount: 30.0,
						},
					}, nil).Times(1)
			},
			wantCode: http.StatusOK,
			want: []shema.HistoryTransfers{
				{
					Time:   "2024-01-23 13:50:19.000000 +00:00",
					From:   "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn",
					To:     "aqaWKX3YeQ6jRWxJ32Zbxwq17kAdEn",
					Amount: 20.0,
				},
				{
					Time:   "2024-01-23 13:50:58.000000 +00:00",
					From:   "ay7uwX3YeQ6jRWxJ32Zbxwq17kAdEn",
					To:     "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn",
					Amount: 30.0,
				},
			},
		},
		{
			name: "BAD1",
			id:   "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAewq",
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GetHistory", mock.Anything, "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAewq").Return(nil, constants.ErrNotFromPerson).Times(1)
			},
			wantCode: http.StatusNotFound,
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()

			service := mocks.NewService(t)
			h := NewHandler(service, config.Config{})

			path := "/api/v1/wallet/:walletId/history"
			g.GET(path, h.History)
			tt.serviceMock(service)
			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/wallet/%s/history", tt.id), nil)

			g.ServeHTTP(w, request)

			if w.Code != tt.wantCode {
				t.Errorf("got %d, want %d", w.Code, tt.wantCode)
			}
			wantResponse, err := json.Marshal(tt.want)
			if err != nil {
				fmt.Errorf("failed json")
				return
			}

			response, err := json.Marshal(w.Body)
			if err != nil {
				fmt.Errorf("failed json")
				return
			}
			if bytes.Equal(wantResponse, response) {
				t.Errorf("got %s, want %v", w.Body, tt.want)
			}
		})
	}
}

func TestHandler_Status(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		serviceMock serviceMock
		wantCode    int
		want        shema.Wallet
	}{
		{
			name: "OK#1",
			id:   "BXbzKE4P62ui6evzzeB5wYLoO1r0Al",
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GetStatus", mock.Anything, "BXbzKE4P62ui6evzzeB5wYLoO1r0Al").Return("BXbzKE4P62ui6evzzeB5wYLoO1r0Al", 120.0, nil).Times(1)
			},
			wantCode: http.StatusOK,
			want:     shema.Wallet{ID: "BXbzKE4P62ui6evzzeB5wYLoO1r0Al", Balance: 120},
		},
		{
			name: "OK#2",
			id:   "BXbzKE4P62ui6qwereB5wYLoO1r0Al",
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GetStatus", mock.Anything, "BXbzKE4P62ui6qwereB5wYLoO1r0Al").Return("BXbzKE4P62ui6qwereB5wYLoO1r0Al", 20.0, nil).Times(1)
			},
			wantCode: http.StatusOK,
			want:     shema.Wallet{ID: "BXbzKE4P62ui6qwereB5wYLoO1r0Al", Balance: 20},
		},
		{
			name: "BAD",
			id:   "BXbzKE4P62ui6qwereB5wYLoO115at",
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GetStatus", mock.Anything, "BXbzKE4P62ui6qwereB5wYLoO115at").Return("", 0.0, constants.ErrNotFromPerson).Times(1)
			},
			wantCode: http.StatusNotFound,
			want:     shema.Wallet{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			service := mocks.NewService(t)
			h := NewHandler(service, config.Config{})
			tt.serviceMock(service)
			path := "/api/v1/wallet/:walletId"
			g.GET(path, h.Status)
			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/wallet/%s", tt.id), nil)

			g.ServeHTTP(w, request)

			if w.Code != tt.wantCode {
				t.Errorf("got %d, want %d", w.Code, tt.wantCode)
			}

			wantResponse, err := json.Marshal(tt.want)
			if err != nil {
				fmt.Errorf("failed json")
				return
			}

			response, err := json.Marshal(w.Body)
			if err != nil {
				fmt.Errorf("failed json")
				return
			}

			if bytes.Equal(wantResponse, response) {
				t.Errorf("got %s, want %v", w.Body, tt.want)
			}
		})
	}
}

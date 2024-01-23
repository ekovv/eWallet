package handler

import (
	"bytes"
	"eWallet/config"
	"eWallet/internal/domains/mocks"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
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
		want        Wallet
	}{
		{
			name: "OK#1",
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GenerateWallet").Return("BXbzKE4P62ui6evzzeB5wYLoO1r0Al", 100.0, nil).Times(1)
			},
			wantCode: http.StatusOK,
			want:     Wallet{ID: "BXbzKE4P62ui6evzzeB5wYLoO1r0Al", Balance: 100},
		},
		{
			name: "OK#2",
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GenerateWallet").Return("BXbzKE4P62ui6qwereB5wYLoO1r0Al", 100.0, nil).Times(1)
			},
			wantCode: http.StatusOK,
			want:     Wallet{ID: "BXbzKE4P62ui6qwereB5wYLoO1r0Al", Balance: 100},
		},
		{
			name: "BAD",
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GenerateWallet").Return("", 0.0, errors.New("can't create id")).Times(1)
			},
			wantCode: http.StatusBadRequest,
			want:     Wallet{},
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
				c.Mock.On("Transaction", "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn", "0MEe1ArlKnX5ea0ByX85PD83QBwpJa", 20.0).Return(nil).Times(1)
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
				c.Mock.On("Transaction", "Oz0WKX3YeQ6jRWxJ32Zbxwq17kAdEn", "jddjkjdksjdkjskd", 20.0).Return(fmt.Errorf("no idOfWallet to person")).Times(1)
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

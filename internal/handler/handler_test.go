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

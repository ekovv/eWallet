package handler

import (
	"eWallet/config"
	"eWallet/internal/domains"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	service domains.Service
	engine  *gin.Engine
	config  config.Config
}

func NewHandler(service domains.Service, cnf config.Config) *Handler {
	router := gin.Default()
	h := &Handler{
		service: service,
		engine:  router,
		config:  cnf,
	}
	Route(router, h)
	return h
}

func (s *Handler) Start() {
	err := s.engine.Run(s.config.Host)
	if err != nil {
		return
	}
}

func (s *Handler) CreateWallet(c *gin.Context) {
	var res Wallet
	id, balance, err := s.service.GenerateWallet()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	res.ID = id
	res.Balance = balance
	bytes, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
	c.Writer.Write(bytes)
}

func (s *Handler) Transactions(c *gin.Context) {
	from := c.Param("walletId")
	var money Transfer
	err := c.ShouldBindJSON(&money)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if money.To == "" || money.Amount == 0.0 {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "bad json"})
		return
	}
	err = s.service.Transaction(from, money.To, money.Amount)
	if err != nil {
		if errors.Is(err, fmt.Errorf("no idOfWallet from person")) {
			c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		if errors.Is(err, fmt.Errorf("no idOfWallet to person")) {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

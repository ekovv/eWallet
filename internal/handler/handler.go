package handler

import (
	"eWallet/config"
	"eWallet/internal/domains"
	"eWallet/internal/shema"
	"encoding/json"
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
	var res shema.Wallet
	id, balance, err := s.service.GenerateWallet()
	if err != nil {
		HandlerErr(c, err)
		return
	}
	res.ID = id
	res.Balance = balance
	bytes, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		HandlerErr(c, err)
		return
	}
	c.Status(http.StatusOK)
	c.Writer.Write(bytes)
}

func (s *Handler) Transactions(c *gin.Context) {
	from := c.Param("walletId")
	var money shema.Transfer
	err := c.ShouldBindJSON(&money)
	if err != nil {
		HandlerErr(c, err)
		return
	}
	err = s.service.Transaction(from, money.To, money.Amount)
	if err != nil {
		HandlerErr(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *Handler) History(c *gin.Context) {
	id := c.Param("walletId")
	history, err := s.service.GetHistory(id)
	if err != nil {
		HandlerErr(c, err)
		return
	}
	var res []shema.HistoryTransfers
	for _, i := range history {
		transfers := shema.HistoryTransfers{}
		transfers.Time = i.Time
		transfers.From = i.From
		transfers.To = i.To
		transfers.Amount = i.Amount
		res = append(res, transfers)
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, res)

}

func (s *Handler) Status(c *gin.Context) {
	id := c.Param("walletId")
	var res shema.Wallet
	idOfWallet, balance, err := s.service.GetStatus(id)
	if err != nil {
		HandlerErr(c, err)
		return
	}
	res.ID = idOfWallet
	res.Balance = balance
	bytes, err := json.MarshalIndent(res, "", "    ")
	c.Status(http.StatusOK)
	c.Writer.Write(bytes)
}

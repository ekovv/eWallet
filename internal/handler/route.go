package handler

import "github.com/gin-gonic/gin"

func Route(c *gin.Engine, h *Handler) {
	c.POST("/api/v1/wallet", h.CreateWallet)
	c.POST("/api/v1/wallet/:walletId/send", h.Transactions)
	c.GET("/api/v1/wallet/:walletId/history", h.History)
}

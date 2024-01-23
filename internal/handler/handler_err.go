package handler

import (
	"eWallet/internal/constants"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandlerErr(c *gin.Context, err error) {
	var UnmarshalTypeError *json.UnmarshalTypeError
	if err != nil {
		if errors.Is(err, constants.ErrNotFromPerson) {
			err := fmt.Sprintf("invalid from person %s", err)
			c.JSON(http.StatusNotFound, err)
			return
		}
		if errors.Is(err, constants.ErrNotToPerson) {
			err := fmt.Sprintf("invalid to person %s", err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		if errors.Is(err, constants.ErrBadAmount) {
			err := fmt.Sprintf("bad amount %s", err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		if errors.As(err, &UnmarshalTypeError) {
			err := fmt.Sprintf("bad json %s", err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusOK)
	return
}

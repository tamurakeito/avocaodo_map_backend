package presentation

import (
	"net/http"

	"github.com/labstack/echo"
)

type HttpHandler struct{}

func NewHttpHandler() HttpHandler {
	handler := HttpHandler{}
	return handler
}

func (handler *HttpHandler) HogeHoge() echo.HandlerFunc {
	type Hoge struct {
		ID int `json:"id"`
	}
	return func(c echo.Context) error {
		models := Hoge{ID: 0}
		return c.JSON(http.StatusOK, models)
	}
}

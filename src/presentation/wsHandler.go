package presentation

import (
	"github.com/labstack/echo"
	"github.com/tamurakeito/avocado_map_backend/src/usecase"
)

type WsHandler struct{}

func NewWsHandler() WsHandler {
	handler := WsHandler{}
	return handler
}

func (handler *WsHandler) WsTextMessage() echo.HandlerFunc {

	return func(c echo.Context) error {

		w := c.Response().Writer
		r := c.Request()

		return usecase.TextMessage(w, r)
	}
}

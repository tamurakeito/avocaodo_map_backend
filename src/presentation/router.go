package presentation

import (
	"net/http"

	"github.com/labstack/echo"
)

func InitRouting(e *echo.Echo, httpHandler HttpHandler, wsHandler WsHandler) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Avocado Map!")
	})
	e.GET("/hogehoge", httpHandler.HogeHoge())

	e.GET("/ws", wsHandler.WsTextMessage())
	e.GET("/location", wsHandler.WsLocation())
}

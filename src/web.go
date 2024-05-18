package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/tamurakeito/avocado_map_backend/src/injector"
	"github.com/tamurakeito/avocado_map_backend/src/presentation"
)

func main() {
	fmt.Println("sever start")
	httpHandler := injector.InjectHttpHandler()
	wsHandler := injector.InjectWsHandler()
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// CORSの設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
	}))

	presentation.InitRouting(e, httpHandler, wsHandler)
	// Logger.Fatalはエラーメッセージをログに出力しアプリケーションを停止する
	// 重要なエラーが発生した場合に使用される
	// 普通のエラーは通常のエラーハンドリングを使おう
	e.Logger.Fatal(e.Start(":8080"))
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"

	//k8s "github.com/SugaoTT/back/k8s"
	handler "github.com/SugaoTT/back/message/handler"
)

func handleWebSocket(c echo.Context) error {
	log.Println("Serving at localhost:8000...")
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		// 初回のメッセージを送信
		err := websocket.Message.Send(ws, "Server: Hello, Next.js!")
		if err != nil {
			c.Logger().Error(err)
		}

		for {
			// Client からのメッセージを読み込む
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			log.Println(msg)
			if err != nil {
				if err.Error() == "EOF" {
					log.Println(fmt.Errorf("read %s", err))
					break
				}
				log.Println(fmt.Errorf("read %s", err))
				c.Logger().Error(err)
			}

			// Client からのメッセージを元に返すメッセージを作成し送信する

			//k8s.Pod_exec(ws)

			handler.HandleMessage(ws, msg)

			/*err = websocket.Message.Send(ws, fmt.Sprintf("Server: \"%s\" received!", msg))
			if err != nil {
				c.Logger().Error(err)
			}*/
			fmt.Println("Server: \"%s\" received!", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/socket", handleWebSocket)

	//外と繋ぐときはlocalhostを外す
	//e.Logger.Fatal(e.Start(":8080"))
	e.Logger.Fatal(e.Start("localhost:8088"))
}

package com

import (
	"eve-firmware/cmds"
	"log"
	"strings"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func InitWS(wg *sync.WaitGroup) {
	cmds.COMMANDS[1].Funcs = append(cmds.COMMANDS[1].Funcs, cmds.CommandFunc{
		NumArgs: 0,
		Func: func(c cmds.CommandCtx) string {
			go StartWS(wg)
			return ""
		},
	})
}

func StartWS(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("WebSocket: \"%s\"", msg)

			res := cmds.ResolveCmds(strings.Split(string(msg), " "))

			if err = c.WriteMessage(mt, []byte(strings.Join(res, "\n"))); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))

	log.Fatal(app.Listen(":3000"))
}

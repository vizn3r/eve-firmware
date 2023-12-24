package com

import (
	"eve-firmware/cmds"
	"log"
	"strings"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func StartWS(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Use("/ws", func (c *fiber.Ctx) error  {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		log.Println(c.Locals("allowed"))  
        log.Println(c.Params("id"))       
        log.Println(c.Query("v"))         
        log.Println(c.Cookies("session"))

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
            log.Printf("recv: %s", msg)

			cmds.ResolveCmds(strings.Split(string(msg), " "))

            if err = c.WriteMessage(mt, msg); err != nil {
                log.Println("write:", err)
                break
            }
        }	
	}))
	
	log.Fatal(app.Listen(":3000"))
}
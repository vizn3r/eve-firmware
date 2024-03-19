package com

import (
	"eve-firmware/arm"
	"eve-firmware/cmds"
	"log"
	"strconv"
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

			res := []string{}
			if strings.HasPrefix(string(msg), "CON") {
				ResolveController(string(msg))
			} else {
				res = cmds.ResolveCmds(strings.Split(string(msg), " "), cmds.FUNCTIONAL)
			}

			if len(res) > 0 {
				if err = c.WriteMessage(mt, []byte(strings.Join(res, "\n"))); err != nil {
					log.Println("write:", err)
					break
				}
			}
		}
	}))

	log.Fatal(app.Listen(":8080"))
}

func ResolveController(data string) {
	rawString := strings.TrimPrefix(data, "CON")
	dataArr := strings.Split(rawString, "/")
	intArr := []int{}
	for _, str := range dataArr {
		num, _ := strconv.Atoi(str)
		intArr = append(intArr, num)
	}

	// Move motors based on controller input
	for i, data := range intArr {
		go arm.MOTORS[i].DoStepDir(arm.MapValue(float64(data), 0, 32768, 100, 0.1), Dir(data))
	}
}

func Dir(data int) int {
	if data > 0 {
		return 1
	}
	return 0
}

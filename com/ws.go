package com

import (
	"eve-firmware/arm"
	"eve-firmware/cmds"
	"eve-firmware/gpio"
	"fmt"
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

			if strings.HasPrefix(string(msg), "CON") {
				ResolveController(string(msg)[3:])
				continue
			}

			res := cmds.ResolveCmds(strings.Split(string(msg), " "), cmds.FUNCTIONAL)
			log.Printf("WebSocket: \"%s\"", msg)
			if err = c.WriteMessage(mt, []byte(strings.Join(res, "\n"))); err != nil {
				fmt.Println(err)
				break
			}
		}
	}))

	log.Fatal(app.Listen(":8080"))
}

func ResolveController(data string) {
	dataArr := strings.Split(data, "/")
	intArr := []int{}
	for _, str := range dataArr {
		num, _ := strconv.Atoi(str)
		intArr = append(intArr, num)
	}

	// Move motors based on controller input
	for i, mot := range arm.MOTORS {
		data := intArr[i]
		if Positive(data) > 10000 && !mot.IsRunning() {
			val := arm.MapValue(float64(Positive(data)), 10000, 32768, 100000, 1)
			if i < 3 {
				val /= 1000
			}
			if err := gpio.Write(mot.Dir, int(Dir(data))); err != nil {
				fmt.Println("Error setting direction:", err)
				return
			}
			go mot.DoStep(val)
		}
	}
}

func Positive(data int) int {
	if data > 0 {
		return data
	}
	return -data
}

func Dir(data int) int {
	if data > 0 {
		return 1
	}
	return 0
}

package visualizer

import (
	"eve-firmware/arm"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Start() {
	app := fiber.New()

	app.Get("/matrix/:id", func(c *fiber.Ctx) error {
		id := c.Params("id", "0")
		m, _ := strconv.Atoi(string(id[0]))
		n, _ := strconv.Atoi(string(id[1]))
		mtx := arm.HTMFromTo(m, n)
		arm.CalculatePosition()
		c.Set("Access-Control-Allow-Origin", "*")
		return c.SendString(mtx.Format())
	})

	app.Get("/script.js", func(c *fiber.Ctx) error {
		return c.SendFile("./visualizer/script.js")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./visualizer/index.html")
	})

	if err := app.Listen(":8080"); err != nil {
		fmt.Println(err)
		return
	}
}

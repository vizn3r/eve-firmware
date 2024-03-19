package com

import (
	"encoding/base64"
	"eve-firmware/cmds"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func StartHTTP() {
	app := fiber.New()

	app.Post("/exec/:file", func(c *fiber.Ctx) error {
		file := c.Params("file")
		out := cmds.ResolveCmds([]string{"s0", "./files/" + file}, cmds.MISC)
		fmt.Println("Executing file '", file, "' from remote client")
		return c.SendString(strings.Join(out, "\n"))
	})

	app.Get("/files", func(c *fiber.Ctx) error {
		dir, err := os.ReadDir("./files")
		if err != nil {
			fmt.Println(err)
			return c.SendString("Error uploading file")
		}
		files := ""
		for _, file := range dir {
			files += file.Name() + " "
		}
		return c.SendString(files)
	})
	app.Get("/files/:file", func(c *fiber.Ctx) error {
		data, _ := os.ReadFile("./files/" + c.Params("file"))
		return c.Send(data)
	})
	app.Post("/files/:file/:data", func(c *fiber.Ctx) error {
		file := c.Params("file")
		data := c.Params("data")

		bytes, _ := base64.URLEncoding.DecodeString(data)
		os.Remove("./files/" + file)
		if err := os.WriteFile("./files/"+file, bytes, os.ModePerm); err != nil {
			return c.SendString("Can't write file")
		}
		fmt.Println("New file'", file, "' received")
		return c.SendStatus(fiber.StatusOK)
	})
	app.Delete("/files/:file", func(c *fiber.Ctx) error {
		file := c.Params("file")
		fmt.Println("Deleting file '", file, "' from remote")
		if err := os.Remove("./files/" + file); err != nil {
			return c.SendString("Can't delete the file")
		}
		return c.SendString("File deleted")
	})
	if err := app.Listen(":8000"); err != nil {
		fmt.Println(err)
	}
}

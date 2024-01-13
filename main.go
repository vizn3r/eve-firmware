package main

import (
	"bufio"
	"eve-firmware/arm"
	"eve-firmware/cmds"
	"eve-firmware/com"
	"eve-firmware/gpio"
	"eve-firmware/util"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const VERSION = "v0.0.3"

const ASCII = "\n" + `_____________    ____________       __________________________ ______  ______       _________ ________ __________
___  ____/__ |  / /___  ____/       ___  ____/____  _/___  __ \___   |/  /__ |     / /___    |___  __ \___  ____/
__  __/   __ | / / __  __/          __  /_     __  /  __  /_/ /__  /|_/ / __ | /| / / __  /| |__  /_/ /__  __/   
_  /___   __ |/ /  _  /___          _  __/    __/ /   _  _, _/ _  /  / /  __ |/ |/ /  _  ___ |_  _, _/ _  /___   
/_____/   _____/   /_____/          /_/       /___/   /_/ |_|  /_/  /_/   ____/|__/   /_/  |_|/_/ |_|  /_____/`

func Clear() {
	switch runtime.GOOS {
	case "windows":
		c := exec.Command("cmd", "/c", "cls")
		c.Stdout = os.Stdout
		if e := c.Run(); e != nil {
			return
		}
	case "linux":
		c := exec.Command("printf", `\033c`)
		c.Stdout = os.Stdout
		if e := c.Run(); e != nil {
			return
		}
	}
}

func main() {
	cmds.COMMANDS = append(cmds.COMMANDS,
		cmds.Command{
			Call: 'T',
			Funcs: []cmds.CommandFunc{
				{
					NumArgs: 1,
					Args:    "<string>",
					Desc:    "Test argument handling",
					Func: func(c cmds.CommandCtx) string {
						return "test return: " + c.Args[0]
					},
				},
				{
					NumArgs: 0,
					Desc:    "Test command handling",
					Func: func(c cmds.CommandCtx) string {
						return "test return 2"
					},
				},
			},
		},
		cmds.Command{
			Call: 'S',
			Funcs: []cmds.CommandFunc{
				{
					NumArgs: 0,
					Desc:    "Exit program",
					Func: func(c cmds.CommandCtx) string {
						if util.Prompt("Do you want to exit?") {
							os.Exit(0)
						}
						return ""
					},
				},
				{
					NumArgs: 0,
					Desc:    "Create service file",
					Func: func(c cmds.CommandCtx) string {
						return ""
					},
				},
			},
		},
		cmds.Command{
			Call: 'H',
			Funcs: []cmds.CommandFunc{
				{
					NumArgs: 0,
					Desc:    "Help menu",
					Func: func(c cmds.CommandCtx) string {
						var out string
						for _, cmd := range cmds.COMMANDS {
							out += string(cmd.Call) + "\n"
							for i, fn := range cmd.Funcs {
								out += "   " + strconv.Itoa(i) + " " + fn.Args + " - " + fn.Desc + "\n"
							}
						}
						return out
					},
				},
			},
		},
		cmds.Command{
			Call: 'G',
			Funcs: []cmds.CommandFunc{
				{
					NumArgs: 2,
					Desc:    "Open pin",
					Args:    "<pin, mode>",
					Func: func(c cmds.CommandCtx) string {
						if err := gpio.Open(c.IntArgs[0]); err != nil {
							return err.Error()
						}
						return ""
					},
				},
				{
					NumArgs: 2,
					Desc:    "Write to pin",
					Args:    "<pin, value>",
					Func: func(c cmds.CommandCtx) string {
						if err := gpio.Write(c.IntArgs[0], c.Args[1]); err != nil {
							return err.Error()
						}
						return ""
					},
				},
			},
		},
	)

	var wg sync.WaitGroup
	go com.InitWS(&wg)

	Clear()

	arm.InitMotors()
	if len(os.Args) > 1 {
		cmds.ResolveCmds(os.Args[1:])
	}

	fmt.Print(ASCII + "\n\nEVE Firmware " + VERSION + "\nby vizn3r 2023\n\n")

	s := bufio.NewScanner(os.Stdin)
	for {
		if s.Scan() {
			cmds.ResolveCmds(strings.Split(strings.TrimSpace(s.Text()), " "))
		}
	}
	wg.Wait()
}

package main

import (
	"bufio"
	"eve-firmware/arm"
	"eve-firmware/cmds"
	"eve-firmware/com"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var ASCII = "\n" + `_____________    ____________       __________________________ ______  ______       _________ ________ __________
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
	var wg sync.WaitGroup
	go com.InitWS(&wg)

	arm.InitMotors()
	if len(os.Args) > 1 {
		cmds.ResolveCmds(os.Args[1:])
	}

	Clear()
	fmt.Print(ASCII + "\n\nEVE Firmware v0.0.2\nby vizn3r 2023\n\n")

	s := bufio.NewScanner(os.Stdin)
	for {
		if s.Scan() {
			cmds.ResolveCmds(strings.Split(strings.TrimSpace(s.Text()), " "))
		}
	}
	wg.Wait()
}

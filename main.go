package main

import (
	"bufio"
	"eve-firmware/arm"
	"eve-firmware/cmds"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var ASCII = "\n" + `_____________    ____________       __________________________ ______  ______       _________ ________ __________
___  ____/__ |  / /___  ____/       ___  ____/____  _/___  __ \___   |/  /__ |     / /___    |___  __ \___  ____/
__  __/   __ | / / __  __/          __  /_     __  /  __  /_/ /__  /|_/ / __ | /| / / __  /| |__  /_/ /__  __/   
_  /___   __ |/ /  _  /___          _  __/    __/ /   _  _, _/ _  /  / /  __ |/ |/ /  _  ___ |_  _, _/ _  /___   
/_____/   _____/   /_____/          /_/       /___/   /_/ |_|  /_/  /_/   ____/|__/   /_/  |_|/_/ |_|  /_____/`


func main() {
	arm.InitMotors()
	if len(os.Args) > 1 {
		cmds.ResolveCmds(os.Args[1:])
	}

	if runtime.GOOS == "windows" {
		c := exec.Command("cmd", "/c", "cls")
		c.Stdout = os.Stdout
		c.Run()
	} else if runtime.GOOS == "linux" {
		c := exec.Command("clear")		
		c.Stdout = os.Stdout
		c.Run()
	}

	fmt.Print(ASCII + "\n\nEVE Firmware v0.0.2\nby vizn3r 2023\n\n")

	var s = bufio.NewScanner(os.Stdin)
	for {
		if s.Scan() {
			cmds.ResolveCmds(strings.Split(strings.TrimSpace(s.Text()), " "))
		}
	}
}
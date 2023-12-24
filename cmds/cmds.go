package cmds

import (
	"eve-firmware/util"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CommandCtx struct {
	Index int
	Args  []string
}

type Command struct {
	Call    byte
	NumArgs []int
	Funcs   []func(c CommandCtx) string
}

/* 	
	DONT FORGET NumArgs WHEN ADDING func
	
	there is no checking for that :P

	should probably rework this in future, but it is fine for now :D
*/
var COMMANDS = []Command {
	{
		// Test commands
		Call: 'T',
		NumArgs: []int{1, 0, 1},		
		Funcs: []func(CommandCtx) string{
			func(c CommandCtx) string {
				return "test return: " + c.Args[0]
			},
			func(c CommandCtx) string {
				return "test return 2"
			},
		},
	},
	{
		Call: 'S',
		NumArgs: []int{0},
		Funcs: []func(c CommandCtx) string{
			func(c CommandCtx) string {
				if util.Prompt("Do you want to exit?") {
					os.Exit(0)		
				}
				return ""
			},
		},
	},
}

func (c *CommandCtx) HasArg(arg string) bool {
	for _, a := range c.Args {
		if strings.EqualFold(arg, a) {
			return true
		}
	}
	return false
}

func CmdHas(cmds []Command, c Command) (bool, Command) {
	for _, cmd := range cmds {
		if cmd.Call == c.Call {
			return true, cmd
		}
	}
	return false, Command{}
}

func ResolveCmds(rawArgs []string) {
	if len(rawArgs) == 0 || rawArgs == nil || rawArgs[0] == "" {
		return
	}
	for i := 0; i < len(rawArgs); i++ {
		a := strings.ToUpper(rawArgs[i])

		var index int
		var cmd Command

		if has, c := CmdHas(COMMANDS, Command{Call: a[0]}); !has {
			fmt.Println("err: Invalid command '" + a + "'")
			continue
		} else { cmd = c }
		if j, e := strconv.Atoi(a[1:]); e != nil || j >= len(cmd.Funcs) {
			fmt.Println("err: Invalid command index '" + a[1:] + "'")
			continue
		} else { index = j }
		if len(rawArgs[i + 1:]) < cmd.NumArgs[index] {
			fmt.Println("err: Not enough args")
			continue
		}

		// []string for Command from rawArgs
		args := rawArgs[i + 1 : i + cmd.NumArgs[index] + 1]

		msg := cmd.Funcs[index](CommandCtx{index, args})
		if msg != "" {
			fmt.Println(a, "out:", msg)
		}

		// Move i by NumArgs to next rawArg
		i += cmd.NumArgs[index]
	}
}
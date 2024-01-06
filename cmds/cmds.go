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

type CommandFunc struct {
	NumArgs int
	Func    func(c CommandCtx) string
}

type Command struct {
	Call  byte
	Funcs []CommandFunc
}

var COMMANDS = []Command{
	{
		// Test commands
		Call: 'T',
		Funcs: []CommandFunc{
			{
				NumArgs: 1,
				Func: func(c CommandCtx) string {
					return "test return: " + c.Args[0]
				},
			},
			{
				NumArgs: 0,
				Func: func(c CommandCtx) string {
					return "test return 2"
				},
			},
		},
	},
	{
		Call: 'S',
		Funcs: []CommandFunc{
			{
				NumArgs: 0,
				Func: func(c CommandCtx) string {
					if util.Prompt("Do you want to exit?") {
						os.Exit(0)
					}
					return ""
				},
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

func ResolveCmds(rawArgs []string) []string {
	if len(rawArgs) == 0 || rawArgs[0] == "" || rawArgs == nil {
		return []string{}
	}
	out := []string{}
	for i := 0; i < len(rawArgs); i++ {
		a := strings.ToUpper(rawArgs[i])

		var index int
		var cmd Command

		if has, c := CmdHas(COMMANDS, Command{Call: a[0]}); !has {
			fmt.Println("err: Invalid command '" + a + "'")
			continue
		} else {
			cmd = c
		}
		if j, e := strconv.Atoi(a[1:]); e != nil || j >= len(cmd.Funcs) {
			fmt.Println("err: Invalid command index '" + a[1:] + "'")
			continue
		} else {
			index = j
		}
		if len(rawArgs[i+1:]) < cmd.Funcs[index].NumArgs {
			fmt.Println("err: Not enough args")
			continue
		}

		// []string for Command from rawArgs
		args := rawArgs[i+1 : i+cmd.Funcs[index].NumArgs+1]

		msg := cmd.Funcs[index].Func(CommandCtx{index, args})
		if msg != "" {
			fmt.Println(msg)
		}

		out = append(out, msg)

		// Move i by NumArgs to next rawArg
		i += cmd.Funcs[index].NumArgs
	}
	return out
}

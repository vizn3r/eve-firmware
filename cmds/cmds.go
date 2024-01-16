package cmds

import (
	"strconv"
	"strings"
)

type CommandCtx struct {
	Index     int
	Args      []string
	IntArgs   []int
	FloatArgs []float64
}

type CommandFunc struct {
	NumArgs int
	Args    string
	Desc    string
	Func    func(c CommandCtx) string
}

type CommandType int

const (
	USER CommandType = iota
	FUNCTIONAL
	MISC
)

type ErrorType string

const (
	INVALID_COMMAND ErrorType = "ERR_INVALID_COMMAND"
	INVALID_INDEX   ErrorType = "ERR_INVALID_INDEX"
	INVALID_ARGS    ErrorType = "ERR_INVALID_ARGS"
)

type Command struct {
	Call  byte
	Funcs []CommandFunc
	Type  CommandType
}

type Commands []Command

var COMMANDS Commands

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

func IsErrorType(s string, t ErrorType) bool {
	err := strings.Split(strings.TrimSpace(s), " ")[0]
	return err == string(t)
}

func IsError(s string) bool {
	err := strings.Split(strings.TrimSpace(s), " ")[0]
	return strings.HasPrefix(err, "ERR")
}

func HasError(s []string) bool {
	for _, i := range s {
		if IsError(i) {
			return true
		}
	}
	return false
}

func ResolveCmds(rawArgs []string, source CommandType) []string {
	if len(rawArgs) == 0 || rawArgs[0] == "" || rawArgs == nil {
		return []string{}
	}
	out := []string{}
	for i := 0; i < len(rawArgs); i++ {
		a := strings.ToUpper(rawArgs[i])

		var index int
		var cmd Command

		if has, c := CmdHas(COMMANDS, Command{Call: a[0]}); !has {
			out = append(out, string(INVALID_COMMAND)+" err: Invalid comand '"+a+"'")
			continue
		} else {
			cmd = c
		}
		if j, e := strconv.Atoi(a[1:]); e != nil || j >= len(cmd.Funcs) {
			out = append(out, string(INVALID_INDEX)+" err: Invalid command index '"+a[:1]+"'")
			continue
		} else {
			index = j
		}
		if len(rawArgs[i+1:]) < cmd.Funcs[index].NumArgs {
			out = append(out, string(INVALID_ARGS)+" err: Not enough args, need "+cmd.Funcs[index].Args)
			continue
		}
		if int(cmd.Type) > int(source) {
			i += cmd.Funcs[index].NumArgs
			continue
		}

		// []string for Command from rawArgs
		args := rawArgs[i+1 : i+cmd.Funcs[index].NumArgs+1]
		var intArgs []int
		var floatArgs []float64
		for _, arg := range args {
			i, err := strconv.Atoi(arg)
			if err != nil {
				i = 0
			}
			j, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				j = 0.0
			}
			intArgs = append(intArgs, i)
			floatArgs = append(floatArgs, j)
		}

		msg := cmd.Funcs[index].Func(CommandCtx{index, args, intArgs, floatArgs})
		out = append(out, msg)

		// Move i by NumArgs to next rawArg
		i += cmd.Funcs[index].NumArgs
	}
	return out
}

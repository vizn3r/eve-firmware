package arm

import (
	"eve-firmware/cmds"
	"eve-firmware/util"
	"strings"
)

// Pins for motor driver
type Motor struct {
	Step int
	Dir int

	Micro1 int
	Micro2 int
	Micro3 int

	Enabe int
}

type MotorConfig struct {
	Motors []Motor
}

var MotorCommands = []cmds.Command{
	
}

// Pull motor config from conf.json and append to COMMANDS
func InitMotors() {
	cmds.COMMANDS = append(cmds.COMMANDS, MotorCommands...)
	cmds.COMMANDS = append(cmds.COMMANDS, cmds.Command{
		Call: 'M',
		Funcs: []cmds.CommandFunc {
			{
				NumArgs: 1,
				Func: func(c cmds.CommandCtx) string {
					switch strings.ToUpper(c.Args[0]) {
					case "F":
						return "Forward"
					case "R":
						return "Reverse"
					}
					return "err: Invalid argument '" + c.Args[0] + "'"
				},
			},
		},
	})
	var motorCfg MotorConfig
	util.ParseJSON("./conf/motors.json", motorCfg)
}
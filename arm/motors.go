package arm

import (
	"eve-firmware/cmds"
	"eve-firmware/gpio"
	"eve-firmware/util"
	"fmt"
	"strconv"
	"time"
)

// Pins for motor driver
type Motor struct {
	Step  int
	Angle float64
	Dir   int
	Diag  int
}

const (
	FORWARD int = iota
	BACKWARD
)

type Motors []Motor

type MotorConfig struct {
	Motors []Motor
}

var MOTORS Motors

var MotorCommands = []cmds.Command{
	{
		Call: 'M',
		Type: cmds.FUNCTIONAL,
		Funcs: []cmds.CommandFunc{
			{
				NumArgs: 0,
				Desc:    "List of all motors",
				Func: func(c cmds.CommandCtx) string {
					var out string
					for i, m := range MOTORS {
						out += strconv.Itoa(i) + ". Motor\n  - STEP: " + strconv.Itoa(m.Step) + "\n  - DIR:  " + strconv.Itoa(m.Dir) + "\n"
					}
					return out
				},
			},

			{
				NumArgs: 4,
				Desc:    "Drive motor by steps",
				Args:    "<motor, dir, steps, delay>",
				Func: func(c cmds.CommandCtx) string {
					MOTORS[c.IntArgs[0]].DriveSteps(c.IntArgs[2], c.FloatArgs[3], int(c.IntArgs[1]))
					return "done"
				},
			},
			{
				NumArgs: 4,
				Desc:    "Drive motor by angle",
				Args:    "<motor, dir, angle, delay>",
				Func: func(c cmds.CommandCtx) string {
					MOTORS[c.IntArgs[0]].DriveAngle(Angle(Angle(c.FloatArgs[2]).Radians()), c.FloatArgs[3], int(c.IntArgs[1]))
					return "done"
				},
			},
			{
				NumArgs: 4,
				Desc:    "Drive motor by steps asynchronously",
				Args:    "<motor, dir, steps, delay>",
				Func: func(c cmds.CommandCtx) string {
					go MOTORS[c.IntArgs[0]].DriveSteps(c.IntArgs[2], c.FloatArgs[3], int(c.IntArgs[1]))
					return "done"
				},
			},
			{
				NumArgs: 4,
				Desc:    "Drive motor by angle asynchronously",
				Args:    "<motor, dir, angle, delay>",
				Func: func(c cmds.CommandCtx) string {
					go MOTORS[c.IntArgs[0]].DriveAngle(Angle(Angle(c.FloatArgs[2]).Radians()), c.FloatArgs[3], int(c.IntArgs[1]))
					return "done"
				},
			},
		},
	},
}

// Load motor configurations from "./conf/motors.json", append them to MOTORS and append MotorCommands to COMMANDS
func InitMotors() {
	cmds.COMMANDS = append(cmds.COMMANDS, MotorCommands...)
	motors := MotorConfig{}
	util.ParseJSON("./conf/motors.json", &motors)
	MOTORS = append(MOTORS, motors.Motors...)
}

// Open GPIO pins of Motor
func (m *Motor) OpenPins() error {
	if err := gpio.Open(m.Step); err != nil {
		return err
	}
	if err := gpio.Open(m.Dir); err != nil {
		return err
	}
	if err := gpio.Write(m.Dir, "1"); err != nil {
		return err
	}
	return nil
}

// Do one step of Motor
func (m *Motor) DoStep(delay float64) error {
	if err := gpio.High(m.Step); err != nil {
		return err
	}
	time.Sleep(time.Millisecond * time.Duration(delay))
	if err := gpio.Low(m.Step); err != nil {
		return err
	}
	time.Sleep(time.Millisecond * time.Duration(delay))
	return nil
}

// Drive motor by Steps with Delay in int
func (m *Motor) DriveSteps(steps int, delay float64, dir int) {
	if err := m.OpenPins(); err != nil {
		fmt.Println(err)
		return
	}
	defer gpio.Close()

	if err := gpio.Write(m.Dir, strconv.Itoa(int(dir))); err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < steps; i++ {
		if err := m.DoStep(delay); err != nil {
			fmt.Println(err)
			return
		}
	}
}

// Drive motor to Angle relative to current position with Delay in int
func (m *Motor) DriveAngle(angle Angle, delay float64, dir int) {
	fmt.Println(angle)
	if err := m.OpenPins(); err != nil {
		fmt.Println(err)
		return
	}
	defer gpio.Close()
	for i := 0.0; i <= angle.Float64(); i += m.Angle {
		if err := m.DoStep(delay); err != nil {
			fmt.Println(err)
			return
		}
	}
}

// Drive motor to Angle relative to current position in Time, in int
func (m *Motor) Drive(angle Angle, time float64, dir int) {
	if err := m.OpenPins(); err != nil {
		fmt.Println(err)
		return
	}
	defer gpio.Close()
	rotations := angle.Degrees() / m.Angle
	delay := time / rotations
	for i := 0.0; i <= angle.Degrees(); i += m.Angle {
		if err := m.DoStep(delay); err != nil {
			fmt.Println(err)
			return
		}
	}
}

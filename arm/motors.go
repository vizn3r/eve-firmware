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
	Step     int
	Angle    float64
	Dir      int
	Stepping SteppingMode

	running bool
}

type MotorConfig struct {
	Motors []Motor
}

type SteppingMode []uint

var (
	// Stepping modes
	SPEED_1 SteppingMode = []uint{0, 0} // 1/8
	SPEED_2 SteppingMode = []uint{1, 1} // 1/16
	SPEED_3 SteppingMode = []uint{0, 1} // 1/32
	SPEED_4 SteppingMode = []uint{1, 0} // 1/64
)

var MOTORS []*Motor

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
				Args:    "<motor, dir, steps, delay [ns]>",
				Func: func(c cmds.CommandCtx) string {
					motor := MOTORS[c.IntArgs[0]]
					motor.DriveSteps(c.IntArgs[2], c.FloatArgs[3], int(c.IntArgs[1]))
					return "motor" + c.Args[0] + " started"
				},
			},
			{
				NumArgs: 4,
				Desc:    "Drive motor by angle",
				Args:    "<motor, dir, angle, delay [ns]>",
				Func: func(c cmds.CommandCtx) string {
					motor := MOTORS[c.IntArgs[0]]
					motor.DriveAngle(Angle(Angle(c.FloatArgs[2]).Radians()), c.FloatArgs[3], int(c.IntArgs[1]))
					return "motor" + c.Args[0] + " started"
				},
			},
			{
				NumArgs: 4,
				Desc:    "Drive motors by steps asynchronously",
				Args:    "<motor, dir, steps, delay [ns]>",
				Func: func(c cmds.CommandCtx) string {
					motor := MOTORS[c.IntArgs[0]]
					for motor.IsRunning() {
						time.Sleep(time.Millisecond)
					}
					go motor.DriveSteps(c.IntArgs[2], c.FloatArgs[3], int(c.IntArgs[1]))
					return "motor" + c.Args[0] + " started in async"
				},
			},
			{
				NumArgs: 4,
				Desc:    "Drive motor by angle asynchronously",
				Args:    "<motor, dir, angle, delay [ns]>",
				Func: func(c cmds.CommandCtx) string {
					go MOTORS[c.IntArgs[0]].DriveAngle(Angle(Angle(c.FloatArgs[2]).Radians()), c.FloatArgs[3], int(c.IntArgs[1]))
					return "motor" + c.Args[0] + " started in async"
				},
			},
		},
	},
}

// Load motor configurations from "./conf/motors.json", append them to MOTORS and append MotorCommands to COMMANDS
func InitMotors() {
	cmds.COMMANDS = append(cmds.COMMANDS, MotorCommands...)
	motors := MotorConfig{}
	util.ParseJSON("./conf/conf.json", &motors)
	for i, motor := range motors.Motors {
		motptr := new(Motor)
		motptr.Step = motor.Step
		motptr.Dir = motor.Dir
		motptr.Angle = motor.Angle
		MOTORS = append(MOTORS, motptr)
		MOTORS[i].running = false
	}
}

// Open GPIO pins of Motor
func (m *Motor) OpenPins() error {
	if err := gpio.Open(m.Step, m.Dir); err != nil {
		return err
	}
	if err := gpio.Write(m.Dir, "1"); err != nil {
		return err
	}
	return nil
}

func (m *Motor) ClosePins() error {
	if err := gpio.Close(m.Step, m.Dir); err != nil {
		return err
	}
	return nil
}

func CloseMotors() {
	for _, motor := range MOTORS {
		if err := motor.ClosePins(); err != nil {
			fmt.Println("Error closing motor pins:", err)
		}
	}
}

func (m *Motor) IsRunning() bool {
	return m.running
}

// Do one step of Motor
func (m *Motor) DoStep(delay float64) error {
	m.running = true
	if err := gpio.High(m.Step); err != nil {
		return err
	}
	time.Sleep(time.Nanosecond * time.Duration(delay))
	if err := gpio.Low(m.Step); err != nil {
		return err
	}
	time.Sleep(time.Nanosecond * time.Duration(delay))
	return nil
}

// Drive motor by Steps with Delay in int
func (m *Motor) DriveSteps(steps int, delay float64, dir int) {
	m.running = true
	if err := m.OpenPins(); err != nil {
		fmt.Println("Error opening pins:", err)
		return
	}

	if err := gpio.Write(m.Dir, int(dir)); err != nil {
		fmt.Println("Error setting direction:", err)
		return
	}

	for i := 0; i < steps; i++ {
		if err := m.DoStep(delay); err != nil {
			fmt.Println("Error doing step:")
		}
	}
	m.running = false
}

// Drive motor to Angle relative to current position with Delay in int
func (m *Motor) DriveAngle(angle Angle, delay float64, dir int) {
	m.running = true
	if err := m.OpenPins(); err != nil {
		fmt.Println("Error opening pins:", err)
		return
	}

	if err := gpio.Write(m.Dir, int(dir)); err != nil {
		fmt.Println("Error doing step:")
	}

	for i := 0.0; i <= angle.Float64(); i += m.Angle {
		if err := m.DoStep(delay); err != nil {
			fmt.Println("Error doing step:")
			return
		}
	}
	m.running = false
}

// Drive motor to Angle relative to current position in Time, in int
func (m *Motor) Drive(angle Angle, time float64, dir int) {
	m.running = true
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
	m.running = false
}

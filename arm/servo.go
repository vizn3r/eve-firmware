package arm

import (
	"eve-firmware/cmds"
	"eve-firmware/gpio"
	"fmt"
	"time"
)

type Servo struct {
	Pin       int
	newPulse  chan float64
	prevPulse float64
	minPulse  float64
	maxPulse  float64

	close chan bool
}

var SERVOS []*Servo

func InitServos() {
	cmds.COMMANDS = append(cmds.COMMANDS, cmds.Command{
		Call: 'E',
		Type: cmds.USER,
		Funcs: []cmds.CommandFunc{
			{
				NumArgs: 2,
				Desc:    "Set servo angle",
				Args:    "<servo, angle>",
				Func: func(c cmds.CommandCtx) string {
					servo := SERVOS[c.IntArgs[0]]
					servo.Angle(c.FloatArgs[1])
					return "Servo angle set"
				},
			},
		},
	})
}

func (s *Servo) Open() {
	s.prevPulse = 1500
	s.minPulse = 1000
	s.maxPulse = 2000
	s.newPulse = make(chan float64)
	s.close = make(chan bool)
	if err := gpio.Open(s.Pin); err != nil {
		fmt.Println("Can't open servo pin:", err)
		return
	}
	if err := gpio.Dir(s.Pin, gpio.OUTPUT); err != nil {
		fmt.Println("Can't set servo pin:", err)
		return
	}
	SERVOS = append(SERVOS, s)
	go func() {
		for {
			if err := gpio.Low(s.Pin); err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Millisecond * 20)
			select {
			case pulse := <-s.newPulse:
				if err := gpio.High(s.Pin); err != nil {
					fmt.Println(err)
				}
				s.prevPulse = pulse
				time.Sleep(time.Microsecond * time.Duration(pulse))
			case <-s.close:
				close(s.close)
				close(s.newPulse)
				return
			default:
				if err := gpio.High(s.Pin); err != nil {
					fmt.Println(err)
				}
				time.Sleep(time.Microsecond * time.Duration(s.prevPulse))
			}
		}
	}()
}

func (s *Servo) Close() {
	s.close <- true
}

func (s *Servo) Angle(angle float64) {
	s.newPulse <- MapValue(angle, 0, 180, s.minPulse, s.maxPulse)
}

func (s *Servo) Pulse(pulse float64) {
	s.newPulse <- pulse
}

func MapValue(value, fromLow, fromHigh, toLow, toHigh float64) float64 {
	return (value-fromLow)*(toHigh-toLow)/(fromHigh-fromLow) + toLow
}

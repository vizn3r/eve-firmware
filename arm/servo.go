package arm

import (
	"os"

	"github.com/stianeikeland/go-rpio/v4"
)

func OpenServo() {
	err := rpio.Open()
	if err != nil {
		os.Exit(1)
	}
	defer rpio.Close()

	pin := rpio.Pin(18)
	pin.Mode(rpio.Pwm)
	pin.Freq(50)
	pin.DutyCycle(15, 100)
}

//	type Servo struct {
//		Pin       int
//		newPulse  chan int
//		prevPulse int
//		minPulse  int
//		maxPulse  int
//
//		close chan bool
//	}
//
// var SERVOS []*Servo
//
//	func InitServos() {
//		cmds.COMMANDS = append(cmds.COMMANDS, cmds.Command{
//			Call: 'E',
//			Type: cmds.USER,
//			Funcs: []cmds.CommandFunc{
//				{
//					NumArgs: 2,
//					Desc:    "Set servo angle",
//					Args:    "<servo, angle>",
//					Func: func(c cmds.CommandCtx) string {
//						servo := SERVOS[c.IntArgs[0]]
//						servo.Angle(c.FloatArgs[1])
//						return "Servo angle set"
//					},
//				},
//			},
//		})
//	}
//
// var (
//
//	pwmChip = "/sys/class/pwm/pwmchip0"
//	pwm0    = "/sys/class/pwm/pwmchip0/pwm0"
//
// )
//
//	func (s *Servo) Open() {
//		if err := gpio.WriteFile(pwmChip+"/export", gpio.Format(0)); err != nil && !strings.HasSuffix(err.Error(), "device or resource busy") {
//			fmt.Println(err)
//			return
//		}
//		if err := gpio.WriteFile(pwm0+"/period", gpio.Format(20000000)); err != nil {
//			fmt.Println(err)
//			return
//		}
//		if err := gpio.WriteFile(pwm0+"/polarity", gpio.Format("normal")); err != nil {
//			fmt.Println(err)
//			return
//		}
//		if err := gpio.WriteFile(pwm0+"/enable", gpio.Format(1)); err != nil {
//			fmt.Println(err)
//			return
//		}
//
//		s.minPulse = 1000000
//		s.maxPulse = 2500000
//		s.prevPulse = s.minPulse
//		s.newPulse = make(chan int)
//		s.close = make(chan bool)
//		if err := gpio.Open(s.Pin); err != nil {
//			fmt.Println("Can't open servo pin:", err)
//			return
//		}
//		if err := gpio.Dir(s.Pin, gpio.OUTPUT); err != nil {
//			fmt.Println("Can't set servo pin:", err)
//			return
//		}
//		SERVOS = append(SERVOS, s)
//		go func() {
//			for {
//				select {
//				case pulse := <-s.newPulse:
//					fmt.Println(pulse)
//					if err := gpio.WriteFile(pwm0+"/duty_cycle", gpio.Format(pulse)); err != nil {
//						fmt.Println(err)
//						return
//					}
//					s.prevPulse = pulse
//				case <-s.close:
//					close(s.close)
//					close(s.newPulse)
//					return
//				default:
//					if err := gpio.WriteFile(pwm0+"/duty_cycle", gpio.Format(s.prevPulse)); err != nil {
//						fmt.Println(err)
//						return
//					}
//				}
//			}
//		}()
//	}
//
//	func (s *Servo) Close() {
//		s.close <- true
//	}
//
//	func (s *Servo) Angle(angle float64) {
//		s.newPulse <- int(MapValue(angle, 0, 180, float64(s.minPulse), float64(s.maxPulse)))
//	}
//
//	func (s *Servo) Pulse(pulse int) {
//		s.newPulse <- pulse
//	}
func MapValue(value, fromLow, fromHigh, toLow, toHigh float64) float64 {
	return (value-fromLow)*(toHigh-toLow)/(fromHigh-fromLow) + toLow
}

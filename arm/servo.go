package arm

import (
	"eve-firmware/cmds"
	"eve-firmware/gpio"
	"fmt"
)

const (
	PWM_PATH  = "/sys/class/pwm/pwmchip0/"
	PWM_PIN   = "0"
	PERIOD    = 20000000 // 20ms
	MIN_PULSE = 500000   // 0.5ms
	MAX_PULSE = 2500000  // 2.5ms
)

var (
	newPulse  chan int
	prevPulse = 0
)

func InitServo() {
	cmds.COMMANDS = append(cmds.COMMANDS, cmds.Command{
		Call: 'E',
		Funcs: []cmds.CommandFunc{
			{
				NumArgs: 1,
				Desc:    "Set servo angle",
				Args:    "<angle>",
				Func: func(c cmds.CommandCtx) string {
					SetAngle(c.FloatArgs[0])
					return gpio.Format("Angle set to ", c.FloatArgs[0])
				},
			},
		},
	})
}

func OpenServo() {
	newPulse = make(chan int)
	OpenPWM(PWM_PIN)

	if err := gpio.WriteFile(PWM_PATH+"pwm"+PWM_PIN+"/polarity", "normal"); err != nil {
		fmt.Println(err)
		return
	}

	if err := gpio.WriteFile(PWM_PATH+"pwm"+PWM_PIN+"/period", fmt.Sprintf("%d", PERIOD)); err != nil {
		fmt.Println(err)
		return
	}

	if err := gpio.WriteFile(PWM_PATH+"pwm"+PWM_PIN+"/enable", "1"); err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			select {
			case pulse := <-newPulse:
				writePulseWidth(pulse)
				prevPulse = pulse
			default:
				writePulseWidth(prevPulse)
			}
		}
	}()
}

func writePulseWidth(width int) {
	if err := gpio.WriteFile(PWM_PATH+"pwm"+PWM_PIN+"/duty_cycle", fmt.Sprintf("%d", width)); err != nil {
		fmt.Println(err)
		return
	}
}

func OpenPWM(pin string) {
	if err := gpio.WriteFile(PWM_PATH+"export", pin); err != nil {
		fmt.Println(err)
		return
	}
}

func SetAngle(angle float64) {
	newPulse <- int(MapValue(angle, 0, 180, MIN_PULSE, MAX_PULSE))
}

func MapValue(value, fromLow, fromHigh, toLow, toHigh float64) float64 {
	return (value-fromLow)*(toHigh-toLow)/(fromHigh-fromLow) + toLow
}

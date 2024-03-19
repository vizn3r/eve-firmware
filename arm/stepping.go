package arm

import (
	"eve-firmware/cmds"
	"eve-firmware/gpio"
	"fmt"
	"time"
)

type Register struct {
	SHCP int // Register clock
	STCP int // Storage clock
	MR   int // Memory reset
	DS   int // Data

	CLKDelay int
}

/* table
SHCP | STCP | MR | Q7S | Qn
---------------------------
 x	 | x 		| L  | L   | NS // Low-level asserted on MR clears shift register. Storage register is unchanged.
 x   | ^    | L  | L   | L  // Empty shift register transferred to storage register.
 ^   | x    | H  | Q6S | NC
 x   | ^    | H  | NC  | QnS
 ^   | ^    | H  | Q6S | QnS

*/

var REGISTERS []*Register

func InitRegisters() {
	cmds.COMMANDS = append(cmds.COMMANDS, cmds.Command{
		Call: 'R',
		Type: cmds.USER,
		Funcs: []cmds.CommandFunc{
			{
				NumArgs: 2,
				Desc:    "Set bits for register",
				Args:    "<register, bits>",
				Func: func(c cmds.CommandCtx) string {
					reg := REGISTERS[c.IntArgs[0]]
					var bits []uint
					for _, strBit := range c.Args[1] {
						switch strBit {
						case '1':
							bits = append(bits, 1)
						case '0':
							bits = append(bits, 0)
						default:
							return "Invalid bit " + string(strBit)
						}
					}
					reg.Write(bits)
					return "bits set"
				},
			},
		},
	})
}

func (r *Register) Open() {
	if err := gpio.Open(r.SHCP, r.STCP, r.MR, r.DS); err != nil {
		fmt.Println("Can't open register:", err)
	}
	REGISTERS = append(REGISTERS, r)
}

func (r *Register) SHCPPulse() {
	if err := gpio.High(r.SHCP); err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Microsecond * time.Duration(r.CLKDelay))
	if err := gpio.Low(r.SHCP); err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Microsecond * time.Duration(r.CLKDelay))
}

func (r *Register) STCPPulse() {
	if err := gpio.High(r.STCP); err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Microsecond * time.Duration(r.CLKDelay))
	if err := gpio.Low(r.STCP); err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Microsecond * time.Duration(r.CLKDelay))
}

func (r *Register) Write(bits []uint) {
	for _, b := range bits {
		if b == 1 {
			if err := gpio.High(r.DS); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := gpio.Low(r.DS); err != nil {
				fmt.Println(err)
			}
		}
		r.SHCPPulse()
	}
	r.STCPPulse()
}

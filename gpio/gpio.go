package gpio

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Mode string

type Pin struct {
	GPIO      int
	Value     any
	Direction Mode
	State     bool
}

const (
	INPUT  Mode = "in"
	OUTPUT Mode = "out"
	NONE   Mode = "none"
	GPIO        = "/sys/class/gpio"
)

var (
	Test   bool
	OFFSET = 0
	pins   = []*Pin{}
)

func NewPin(gpio int) *Pin {
	if pin := FindPin(gpio); pin != nil {
		return pin
	}
	pin := new(Pin)
	pin.GPIO = gpio + OFFSET
	pin.Direction = NONE
	pin.State = true
	pins = append(pins, pin)
	return pin
}

func FindPin(gpio int) *Pin {
	for _, pin := range pins {
		if pin.GPIO == gpio+OFFSET {
			return pin
		}
	}
	return nil
}

func ClosePin(gpio int) {
	pin := FindPin(gpio)
	err := WriteFile(GPIO+"/unexport", Format(pin.GPIO))
	if err != nil {
		fmt.Println(err)
		return
	}
	pin.State = false
}

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func WriteFile(file string, args string) error {
	if Test {
		fmt.Println("Writing", args, "to", file)
		return nil
	}
	if len(args) < 1 {
		return fmt.Errorf("not enough arguments")
	}
	err := os.WriteFile(file, []byte(args), 0777)
	if os.IsPermission(err) {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

func ReadFile(file string) ([]byte, error) {
	if Test {
		fmt.Println("Reading file", file)
		return []byte("test"), nil
	}
	mem, err := os.OpenFile(file, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer mem.Close()
	r := bufio.NewReader(mem)
	buf := make([]byte, r.Size())
	if _, err := r.Read(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func Format(args ...any) string {
	out := ""
	for _, arg := range args {
		if reflect.TypeOf(arg) == reflect.TypeFor[int]() {
			out += strconv.Itoa(arg.(int))
		} else if reflect.TypeOf(arg) == reflect.TypeFor[float64]() {
			out += strconv.FormatFloat(arg.(float64), 'G', 5, 64)
		} else {
			out += arg.(string)
		}
	}
	return out
}

func Open(gpios ...int) error {
	// Find offset
	if Test {
		fmt.Println("Reading dir", GPIO)
	} else {
		files, err := os.ReadDir(GPIO)
		if err != nil {
			return err
		}
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "gpiochip") {
				raw := strings.TrimPrefix(file.Name(), "gpiochip")
				num, err := strconv.Atoi(raw)
				if err != nil {
					return err
				}
				OFFSET = num
			}
		}
	}
	for _, pin := range gpios {
		NewPin(pin)
		if err := WriteFile(Format(GPIO, "/export"), Format(FindPin(pin).GPIO)); err != nil && !strings.HasSuffix(err.Error(), "device or resource busy") {
			return err
		}
	}
	return nil
}

func Close(gpios ...int) error {
	for _, pin := range gpios {
		if err := WriteFile(GPIO+"/unexport", Format(FindPin(pin).GPIO)); err != nil {
			return err
		}
		FindPin(pin).State = false
	}
	return nil
}

func CloseAll() error {
	for _, pin := range pins {
		ClosePin(pin.GPIO)
	}
	return nil
}

func Dir(pin int, mode Mode) error {
	if err := WriteFile(Format(GPIO, "/gpio", FindPin(pin).GPIO, "/direction"), string(mode)); err != nil {
		return err
	}
	FindPin(pin).Direction = mode
	return nil
}

func Write(pin int, value any) error {
	if FindPin(pin).Direction != OUTPUT {
		if err := Dir(pin, OUTPUT); err != nil {
			return err
		}
	}
	if err := WriteFile(Format(GPIO, "/gpio", FindPin(pin).GPIO, "/value"), Format(value)); err != nil {
		return err
	}
	FindPin(pin).Value = value
	return nil
}

func High(pin int) error {
	if err := Write(pin, 1); err != nil {
		return err
	}
	return nil
}

func Low(pin int) error {
	if err := Write(pin, 0); err != nil {
		return err
	}
	return nil
}

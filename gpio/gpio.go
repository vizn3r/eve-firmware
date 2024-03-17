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
	pins   = make(map[int]*Pin)
)

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func WriteFile(file string, args string) error {
	if Test {
		fmt.Println("Writing to", file)
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

func format(args ...any) string {
	out := ""
	for _, arg := range args {
		if reflect.TypeOf(arg) == reflect.TypeFor[int]() {
			out += strconv.Itoa(arg.(int))
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
		if err := WriteFile(format(GPIO, "/export"), format(OFFSET+pin)); err != nil && !strings.HasSuffix(err.Error(), "device or resource busy") {
			return err
		}
		pins[OFFSET+pin] = new(Pin)
		pins[OFFSET+pin].State = true
		pins[OFFSET+pin].Value = -1
		pins[OFFSET+pin].Direction = NONE
	}
	return nil
}

func Close(gpios ...int) error {
	for _, pin := range gpios {
		if err := WriteFile(GPIO+"/unexport", format(OFFSET+pin)); err != nil {
			return err
		}
		pins[OFFSET+pin].State = false
	}
	return nil
}

func CloseAll() error {
	for gpio, pin := range pins {
		err := WriteFile(GPIO+"/unexport", format(OFFSET+gpio))
		if err != nil {
			return err
		}
		pin.State = false
	}
	return nil
}

func Dir(pin int, mode Mode) error {
	if err := WriteFile(format(GPIO, "/gpio", OFFSET+pin, "/direction"), string(mode)); err != nil {
		return err
	}
	pins[OFFSET+pin].Direction = mode
	return nil
}

func Write(pin int, value any) error {
	if pins[OFFSET+pin].Direction != OUTPUT {
		if err := Dir(pin, OUTPUT); err != nil {
			return err
		}
	}
	if err := WriteFile(format(GPIO, "/gpio", OFFSET+pin, "/value"), format(value)); err != nil {
		return err
	}
	pins[OFFSET+pin].Value = value
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

// func High(pin int) error {
// 	if err := Write(pin, "1"); err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func Low(pin int) error {
// 	if err := Write(pin, "0"); err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func Write(pin int, value string) error {
// 	if err := Dir(pin, OUTPUT); err != nil {
// 		return err
// 	}
// 	err := WriteFile(GPIO+"/gpio"+fmt.Sprintf("%d", OFFSET+pin)+"/value", value)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func Read(pin int) (string, error) {
// 	if err := Dir(pin, INPUT); err != nil {
// 		return "", err
// 	}
// 	b, err := ReadFile(GPIO + "/gpio" + strconv.Itoa(OFFSET+pin) + "/value")
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(b), nil
// }

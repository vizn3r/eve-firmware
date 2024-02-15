package gpio

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Mode string

const (
	INPUT  Mode = "in"
	OUTPUT Mode = "out"
	GPIO        = "/sys/class/gpio"
)

var (
	open []int
	Test bool
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

func Dir(pin int, mode Mode) error {
	if err := WriteFile(GPIO+"/gpio"+strconv.Itoa(pin)+"/direction", string(mode)); err != nil {
		return err
	}
	return nil
}

func Open(pin int) error {
	err := WriteFile(GPIO+"/export", fmt.Sprintf("%d", pin))
	if err != nil && !strings.HasSuffix(err.Error(), "device or resource busy") {
		return err
	}
	open = append(open, pin)
	return nil
}

func Close() error {
	for _, pin := range open {
		err := WriteFile(GPIO+"/unexport", fmt.Sprintf("%d", pin))
		if err != nil {
			open = []int{}
			return err
		}
	}
	open = []int{}
	return nil
}

func High(pin int) error {
	if err := Write(pin, "1"); err != nil {
		return err
	}
	return nil
}

func Low(pin int) error {
	if err := Write(pin, "0"); err != nil {
		return err
	}
	return nil
}

func Write(pin int, value string) error {
	if err := Dir(pin, OUTPUT); err != nil {
		return err
	}
	err := WriteFile(GPIO+"/gpio"+fmt.Sprintf("%d", pin)+"/value", value)
	if err != nil {
		return err
	}
	return nil
}

func Read(pin int) (string, error) {
	if err := Dir(pin, INPUT); err != nil {
		return "", err
	}
	b, err := ReadFile(GPIO + "/gpio" + strconv.Itoa(pin) + "/value")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

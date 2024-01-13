package util

import (
	"encoding/json"
	"fmt"
	"os"
)

func ParseJSON(file string, data any) {
	byte, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("err: Can't open JSON file", file, err)
		return
	}
	err = json.Unmarshal(byte, &data)
	if err != nil {
		fmt.Println("err: Can't unmarshall JSON file", file, err)
		return
	}
}

func ToJSON(file string, data any) {
	byte, err := json.Marshal(data)
	if err != nil {
		fmt.Println("err: Can't marshall JSON data", data, err)
		return
	}
	err = os.WriteFile(file, byte, 0777)
	if err != nil {
		fmt.Println("err: Can't write JSON file", file, err)
	}
}


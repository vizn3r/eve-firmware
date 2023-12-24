package util

import (
	"fmt"
	"strings"
)

func Prompt(msg string) bool {
	fmt.Print(msg + "[y / n]: ")
	var ans string
	fmt.Scanln(&ans)
	return strings.EqualFold(ans, msg) || strings.HasPrefix(strings.ToLower(ans), "y")
}	
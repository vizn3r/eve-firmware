package util

import (
	"bufio"
	"eve-firmware/cmds"
	"os"
	"strconv"
	"strings"
	"time"
)

// func EveCommands() {
// 	cmds.COMMANDS = append(cmds.COMMANDS, cmds.Command{
// 		Call: 'E',
// 		Type: cmds.MISC,
// 		Funcs: []cmds.CommandFunc{
// 			{
// 				NumArgs: 1,
// 				Func: func(c cmds.CommandCtx) string {
// 					return "jump " + c.Args[0]
// 				},
// 			},
// 			{
// 				NumArgs: 1,
// 				Func: func(c cmds.CommandCtx) string {
// 					return "delay " + c.Args[0]
// 				},
// 			},
// 		},
// 	})
// }

func getInt(s string) int {
	for _, j := range strings.Split(strings.TrimSpace(s), " ") {
		k, err := strconv.Atoi(j)
		if err == nil {
			return k
		}
	}
	return 0
}

func getString(s string, i int) string {
	return strings.Split(strings.TrimSpace(s), " ")[i]
}

func hasVar(s string) bool {
	for _, v := range strings.Split(strings.TrimSpace(s), " ") {
		if strings.HasPrefix(v, "$") {
			return true
		}
	}
	return false
}

func resolveVar(s string, vars map[string]string) []string {
	out := []string{}
	for _, v := range strings.Split(strings.TrimSpace(s), " ") {
		if strings.HasPrefix(v, "$") {
			out = append(out, vars[v[1:]])
		} else {
			out = append(out, v)
		}
	}
	return out
}

func resolveIf(s string, vars map[string]string) bool {
	res := resolveVar(s, vars)
	switch res[1] {
	case "==":
		return res[0] == res[2]
	case "!=":
		return res[0] != res[2]
	}
	return false
}

func jumpEndIf(b []string, start int) int {
	for i := start; i < len(b); i++ {
		if strings.ToLower(b[i]) == "endif" {
			return i
		}
	}
	return len(b)
}

func EveDecode(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return err.Error()
	}
	defer f.Close()

	var buff []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		buff = append(buff, s.Text())
	}
	outBuff := ""

	varBuff := make(map[string]string)
	for i := 0; i < len(buff); i++ {
		switch strings.ToLower(getString(buff[i], 0)) {
		case "print":
			outBuff += strings.Join(resolveVar(buff[i][len("print")+2:], varBuff), " ")
		case "println":
			outBuff += strings.Join(resolveVar(buff[i][len("printf")+2:], varBuff), " ") + "\n"
		case "jump":
			i = getInt(buff[i]) - 2
		case "sleep":
			time.Sleep(time.Millisecond * time.Duration(getInt(buff[i])))
		case "var":
			switch getString(buff[i], 1) {
			case "++":
				n, _ := strconv.Atoi(varBuff[getString(buff[i], 2)])
				n++
				varBuff[getString(buff[i], 2)] = strconv.Itoa(n)
			case "--":
				n, _ := strconv.Atoi(varBuff[getString(buff[i], 2)])
				n--
				varBuff[getString(buff[i], 2)] = strconv.Itoa(n)
			default:
				varBuff[getString(buff[i], 1)] = strings.Join(strings.Split(strings.TrimSpace(buff[i]), " ")[2:], " ")
			}
		case "if":
			if !resolveIf(strings.Join(strings.Split(strings.TrimSpace(buff[i]), " ")[1:], " "), varBuff) {
				i = jumpEndIf(buff, i)
			}
		case "endif":
		default:
			out := cmds.ResolveCmds(resolveVar(buff[i], varBuff), cmds.MISC)
			if !cmds.HasError(out) {
				outBuff += strings.Join(out, "\n")
			}
		}
	}
	return outBuff
}

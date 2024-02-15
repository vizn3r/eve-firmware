package log

import (
	"eve-firmware/util"
	"fmt"
	"reflect"
)

type Log struct {
	Mode LogMode
}

type LogMode int

const (
	MODE_NORMAL LogMode = iota
	MODE_VERBOSE
	MODE_SILENT
	MODE_DEBUG
)

type LogMessage struct {
	Type    LogType
	Normal  []any
	Verbose []any
	Error   LogError
}

type LogError struct {
	Message []any
	Code    int
}

type LogType int

const (
	TYPE_MESSAGE LogType = iota
	TYPE_ERROR
	TYPE_INFO
	TYPE_LOADING
	TYPE_WARNING
	TYPE_DEBUG
)

var LOG = new(Log)

func (l LogMessage) logMode() []any {
	switch LOG.Mode {
	case MODE_SILENT:
		return nil
	case MODE_NORMAL:
		return l.Normal
	case MODE_VERBOSE:
		if l.Verbose != nil {
			return l.Verbose
		}
		return l.Normal
	}
	return nil
}

func InitLog() {
	var LogConfig Log
	util.ParseJSON("./conf/log.json", &LogConfig)
	LOG = &LogConfig
}

func Print(args ...any) {
	for _, a := range args {
		if reflect.TypeOf(a).Kind() == reflect.String && LOG.Mode != MODE_SILENT {
			fmt.Print(a)
		} else {
			l := a.(LogMessage)
			switch l.Type {
			case TYPE_MESSAGE:
				fmt.Print("MSG:", l.logMode()...)
			case TYPE_INFO:
				fmt.Print("INFO:", l.logMode()...)
			case TYPE_ERROR:
				fmt.Print("ERR", l.Error.Code, l.Error.Message...)
			}
		}
	}
}

func Println(args ...any) {
	Print(args...)
	fmt.Print("\n")
}

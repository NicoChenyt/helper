package helper

import (
	"fmt"
	"time"
	"os"
	"strings"
)


var ShowDebug = false


func StdLog(args ... interface{}) (err error) {
	var level string
	var logContent []interface{}
	if len(args) == 1 {
		logContent = []interface{}{args[0]}
		level = "info"
	} else {
		last := args[len(args)-1]
		if last == "info" || last == "error" || last == "warning" || last == "debug" {
			logContent = args[:len(args)-1]
			level = last.(string)
		} else {
			logContent = args
			level = "error"
		}
	}

	host, e := os.Hostname()
	if e != nil {
		host = "*Unknown*"
	}

	fmt.Printf("[%s] (from:%s) %s: ", time.Now().Format("2006-01-02T15:04:05"), host, strings.ToUpper(level))
	fmt.Println(logContent...)

	//if level == "error" {
	//	fmt.Fprintln(os.Stderr, preLog, logContent)
	//}

	return
}

func StdDebug(args ... interface{}) (err error) {
	if ShowDebug {
		if args[len(args) - 1] != "debug" {
			args = append(args, "debug")
		}
		return StdLog(args...)
	}
	return
}

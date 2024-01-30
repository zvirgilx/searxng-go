package util

import (
	"fmt"
	"runtime"
	"time"
)

// catchPanicMessage collect the message of error and stack of program panic.
func catchPanicMessage(err any) any {
	const size = 64 << 10
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]
	err = fmt.Sprintf("%v\n%s", err, buf)
	return err
}

// RecoverFromPanic Capture the panic of the program and print the stack information to the console.
func RecoverFromPanic() {
	if err := recover(); err != nil {
		fmt.Printf("%s[Recovery] %s panic recovered:\n%s\n%s",
			"\033[31m", time.Now().Format("2006/01/02 - 15:04:05"), catchPanicMessage(err), "\033[0m")
	}
}

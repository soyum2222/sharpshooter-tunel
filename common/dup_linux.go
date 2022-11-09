package common

import (
	"os"
	"syscall"
)

func Dup() {
	logFile, err := os.OpenFile("/tmp/fatal.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		panic(err)
		return
	}
	syscall.Dup2(int(logFile.Fd()), int(os.Stderr.Fd()))
}

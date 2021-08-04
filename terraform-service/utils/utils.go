package utils

import (
	"log"
	"os"
	"time"
)

func Log(msg string) {
	l := log.New(os.Stdout, "", 0)
	l.SetPrefix(time.Now().Format("2006-01-02 15:04:05") + " ")
	l.Print(msg)
}

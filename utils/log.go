package utils

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logs = logrus.New()

func init() {
	Logs.Out = os.Stdout

	file, err := os.OpenFile("log/log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Logs.Out = file
	} else {
		panic("Failed to log to file")
	}
}

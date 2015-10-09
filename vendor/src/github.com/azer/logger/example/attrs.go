package main

import (
	"errors"
	"time"

	"github.com/azer/logger"
)

var log = logger.New("e-mail")

func main() {
	log.Info("Sending an e-mail", logger.Attrs{
		"from": "foo@bar.com",
		"to":   "qux@corge.com",
	})

	err := errors.New("Too busy")

	log.Error("Failed to send e-mail. Error: %s", err, logger.Attrs{
		"from": "foo@bar.com",
		"to":   "qux@corge.com",
	})

	timer := log.Timer()
	time.Sleep(time.Millisecond * 500)
	timer.End("Created a new %s image", "bike", logger.Attrs{
		"id":    123456,
		"model": "bmx",
		"frame": "purple",
		"year":  2014,
	})
}

package main

import (
	"time"

	"github.com/azer/logger"
)

var log = logger.New("app")

func main() {
	log.Info("Starting at %d", 9088)

	log.Info("Requesting an image at foo/bar.jpg")
	timer := log.Timer()
	time.Sleep(time.Millisecond * 250)
	timer.End("Fetched foo/bar.jpg")

	log.Error("Failed to start, shutting down...")
}

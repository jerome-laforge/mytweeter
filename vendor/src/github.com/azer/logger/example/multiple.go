package main

import (
	"errors"
	"time"

	"github.com/azer/logger"
)

var app = logger.New("app")
var images = logger.New("images")
var socket = logger.New("websocket")
var users = logger.New("users")
var db = logger.New("database")

func main() {
	app.Info("Starting at %d", 9088)
	db.Info("Connecting to mysql://azer@localhost:9900/foobar")

	images.Info("Requesting an image at foo/bar.jpg")
	timer := images.Timer()
	time.Sleep(time.Millisecond * 250)
	timer.End("Fetched foo/bar.jpg")

	db.Error("Fatal connection error.")

	users.Info("%s just logged  from %s", "John", "Istanbul")

	socket.Info("Connecting...")

	err := errors.New("Unable to connect.")
	socket.Error("%v", err)

	time.Sleep(time.Millisecond * 250)

	app.Error("Failed to start, shutting down...")
}

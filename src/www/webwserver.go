package www

import "sync"

type WebServer struct {
}

var (
	instance *WebServer
	once     sync.Once
)

func GetWebServer() {
	once.Do(func() {
		instance = new(WebServer)

	})
}

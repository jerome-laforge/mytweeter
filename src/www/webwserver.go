package www

import (
	"dto"
	"net/http"
	"sync"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type WebServer struct {
	e *echo.Echo
}

var (
	instance *WebServer
	once     sync.Once
)

func GetWebServer() {
	once.Do(func() {
		instance = new(WebServer)
		instance.e = echo.New()
		instance.e.Use(middleware.Logger())
		instance.e.Get("/list/:id", getAllTweetFor)
		instance.e.Static("/", "www/static")
		instance.e.Run(":8080")
	})
}

func getAllTweetFor(c *echo.Context) error {
	id := c.Param("id")
	tw := dto.NewTweet(id, "test message json")
	return c.JSON(http.StatusOK, tw)
}

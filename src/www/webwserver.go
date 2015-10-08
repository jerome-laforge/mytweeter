package www

import (
	"dto"
	"log"
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
		instance.e.Post("/api/v1/tweet", createTweetV1)
		instance.e.Get("/api/v1/tweets/:id", getAllTweetForV1)
		//instance.e.Static("/", "/www/static")
		instance.e.Run(":8080")
	})
}

func getAllTweetForV1(c *echo.Context) error {
	id := c.Param("id")
	tweets, err := dto.GetAllTweetsForTimeLine(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if len(tweets) == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, tweets)
}

func createTweetV1(c *echo.Context) error {
	tweet := new(dto.Tweet)
	err := c.Bind(tweet)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusBadRequest, nil)
	}
	tweet.GenerateId()
	tweet.Insert()
	return c.JSON(http.StatusOK, tweet)
}

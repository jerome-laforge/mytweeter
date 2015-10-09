package www

import (
	"config"
	"dto"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func StartWebServer() error {
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	e.Post("/api/v1/tweet", createTweetV1)
	e.Get("/api/v1/tweets/:id", getAllTweetForV1)
	//e.Static("/", "/www/static")
	e.Run(conf.Web.Address)
	return nil
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

package www

import (
	"config"
	"dto"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/inconshreveable/log15"
	"github.com/labstack/echo"
)

var log log15.Logger

func StartWebServer() error {
	log = log15.New("module", "webserver")
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}
	e := echo.New()
	e.Post("/api/v1/tweet", createTweetV1)
	e.Get("/api/v1/tweets/:id", getAllTweetForV1)
	e.Get("/api/v1/wait/:timeout", waitFor)
	//e.Static("/", "/www/static")
	log.Info("Launching server on " + conf.Web.Address)
	err = http.ListenAndServe(conf.Web.Address, handlers.LoggingHandler(os.Stdout, handlers.CompressHandler(e.Router())))
	if err != nil {
		log.Error(fmt.Sprintf("Error during start web server : %s", err))
	}
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
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}
	tweet.GenerateId()
	tweet.Insert()
	return c.JSON(http.StatusOK, tweet)
}

func waitFor(c *echo.Context) error {
	timeout, err := time.ParseDuration(c.Param("timeout"))
	if err != nil {
		timeout = 500 * time.Millisecond
	}

	time.Sleep(timeout)
	return c.JSON(http.StatusOK, timeout.String())
}

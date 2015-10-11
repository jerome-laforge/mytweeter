package www

import (
	"config"
	"dto"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/fvbock/endless"
	"github.com/gorilla/handlers"
	"github.com/inconshreveable/log15"
	"github.com/labstack/echo"
)

var (
	log        log15.Logger
	terminated int32
)

func StartWebServer() error {
	log = log15.New("module", "webserver")
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}

	var hystrixTimeout time.Duration
	conf.Hystrix.Timeout = strings.TrimSpace(conf.Hystrix.Timeout)
	if conf.Hystrix.Timeout != "" {
		hystrixTimeout, err = time.ParseDuration(conf.Hystrix.Timeout)
		if err != nil || hystrixTimeout < time.Millisecond {
			hystrixTimeout = time.Second
			log.Error(fmt.Sprintf("Use default time for hystrix timeout %s", hystrixTimeout))
		}
	}

	hystrix.ConfigureCommand("waitFor", hystrix.CommandConfig{
		Timeout:                int(int64(hystrixTimeout) / int64(time.Millisecond)), // converted into Millisecond.
		MaxConcurrentRequests:  conf.Hystrix.MaxConcurrentRequests,
		ErrorPercentThreshold:  conf.Hystrix.ErrorPercentThreshold,
		RequestVolumeThreshold: conf.Hystrix.RequestVolumeThreshold,
		SleepWindow:            conf.Hystrix.SleepWindow,
	})

	conf.Endless.DefaultHammerTime = strings.TrimSpace(conf.Endless.DefaultHammerTime)
	if conf.Endless.DefaultHammerTime != "" {
		duration, err := time.ParseDuration(conf.Endless.DefaultHammerTime)
		if err == nil {
			endless.DefaultHammerTime = duration
		} else {
			log.Error("Bad format for Endless/DefaultHammerTime " + conf.Endless.DefaultHammerTime + " err: " + err.Error())
		}
	}

	e := echo.New()
	e.Post("/api/v1/tweet", createTweetV1)
	e.Get("/api/v1/tweets/:id", getAllTweetForV1)
	e.Get("/api/v1/wait/:timeout", waitFor)
	e.Get("/api/v1/wait_protected/:timeout", waitForProtected)
	//e.Static("/", "/www/static")
	log.Info(fmt.Sprintf("Launching server [pid=%s] on %s", strconv.Itoa(os.Getpid()), conf.Web.Address))
	//err = endless.ListenAndServe(conf.Web.Address, handlers.LoggingHandler(os.Stdout, handlers.CompressHandler(e.Router())))
	srv := endless.NewServer(conf.Web.Address, handlers.LoggingHandler(os.Stdout, handlers.CompressHandler(e.Router())))
	preHookFunc := func() {
		atomic.StoreInt32(&terminated, 1)
	}
	srv.RegisterSignalHook(endless.PRE_SIGNAL, syscall.SIGHUP, preHookFunc)
	srv.RegisterSignalHook(endless.PRE_SIGNAL, syscall.SIGINT, preHookFunc)
	srv.RegisterSignalHook(endless.PRE_SIGNAL, syscall.SIGTERM, preHookFunc)
	err = srv.ListenAndServe()
	if atomic.LoadInt32(&terminated) == 0 {
		if err != nil {
			log.Error(fmt.Sprintf("During startup of server [pid=%s], this error has occurred : %s", strconv.Itoa(os.Getpid()), err))
		}
		return err
	} else {
		log.Info(fmt.Sprintf("Server [pid=%s] is going to shutdown", strconv.Itoa(os.Getpid())))
		return nil
	}
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
	err = tweet.Insert()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, tweet)
}

func waitFor(c *echo.Context) error {
	timeout, err := time.ParseDuration(c.Param("timeout"))
	if err != nil {
		timeout = 500 * time.Millisecond
	}

	time.Sleep(timeout)
	return c.JSON(http.StatusOK, timeout.String()+" by pid = "+strconv.Itoa(os.Getpid()))
}

func waitForProtected(c *echo.Context) error {
	var response *http.Response
	hystrix.Do("waitFor", func() error {
		var err error
		response, err = http.Get("http://127.0.0.1:8080/api/v1/wait/" + c.Param("timeout"))
		//response, err = http.Get(fmt.Sprintf("%s://%s%s", c.Request().URL.Scheme, c.Request().URL.Host, c.Request().URL.Path))
		if err != nil {
			return err
		}
		r := response.Body
		w := c.Response().Writer()
		io.Copy(w, r)
		return nil
	}, func(err error) error {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return nil
	})

	return nil
}

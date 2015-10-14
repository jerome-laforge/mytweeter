package www

import (
	"config"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/fvbock/endless"
)

func listenAndServer(addr string, handler http.Handler) error {
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}

	conf.Endless.DefaultHammerTime = strings.TrimSpace(conf.Endless.DefaultHammerTime)

	if conf.Endless.DefaultHammerTime != "" {
		duration, err := time.ParseDuration(conf.Endless.DefaultHammerTime)
		if err == nil {
			endless.DefaultHammerTime = duration
		} else {
			log.Error("Bad format for Endless/DefaultHammerTime " + conf.Endless.DefaultHammerTime + " err: " + err.Error())
		}
	}

	var terminated int32
	srv := endless.NewServer(addr, handler)
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

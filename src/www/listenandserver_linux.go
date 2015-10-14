package www

import (
	"config"
	"net/http"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/fvbock/endless"
	"github.com/inconshreveable/log15"
)

func listenAndServer(log log15.Logger, addr string, handler http.Handler) error {
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
			log.Error("Bad format", log15.Ctx{"module": "Endless", "DefaultHammerTime": conf.Endless.DefaultHammerTime, "error": err})
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

	log.Info("Launching server")
	err = srv.ListenAndServe()
	if atomic.LoadInt32(&terminated) == 0 {
		if err != nil {
			log.Error("During startup, error has occurred", "error", err)
		}
		return err
	} else {
		log.Info("Server is going to shutdown")
		return nil
	}
}

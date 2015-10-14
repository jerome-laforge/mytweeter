package www

import (
	"net/http"

	"github.com/inconshreveable/log15"
)

func listenAndServer(log log15.Logger, addr string, handler http.Handler) error {
	log.Info("Launching server")
	err := http.ListenAndServe(addr, handler)
	log.Error("During startup, error has occurred", "error", err)
	return err
}

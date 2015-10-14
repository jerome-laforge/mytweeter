package www

import "net/http"

func listenAndServer(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}

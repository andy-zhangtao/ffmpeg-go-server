package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"ffmpeg-go-server/web/avfilter"
	"github.com/gorilla/mux"
)

const version = "/v1/"

func FFServer(port int) (*http.Server) {

	funcMap := register()

	router := mux.NewRouter()
	router.HandleFunc("/_ping", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler

		var ps []string
		for path := range funcMap {
			ps = append(ps, path)
		}

		json.NewEncoder(w).Encode(ps)
	})

	for path, f := range register() {
		router.HandleFunc(path, f)
	}

	addr := fmt.Sprintf("%s:%d", "0.0.0.0", port)

	srv := &http.Server{
		Handler: router,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv
}

func register() (map[string]func(http.ResponseWriter, *http.Request)) {

	funcMap := make(map[string]func(http.ResponseWriter, *http.Request))
	funcMap[wrapAPI("health")] = healthFunc

	registerAVFilter(avfilter.Ifade{}, &funcMap)

	return funcMap
}

func healthFunc(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"name": "healthFunc"})
}

func wrapAPI(path string) string {
	return fmt.Sprintf("%s%s", version, path)
}

func registerAVFilter(i IFilter, fmPtr *map[string]func(http.ResponseWriter, *http.Request)) {
	p, f := i.RegisterInfo()
	pp := wrapAPI(p)
	(*fmPtr)[pp] = f
}

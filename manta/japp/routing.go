package japp

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// InitRouter router
func InitRouter() *mux.Router {
	router := mux.NewRouter()
	dir, _ := os.Getwd()

	buildDir := dir + "/manta/frontend/dist/"
	static := buildDir + "static/"

	// Index file
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, buildDir+"/index.html")
	})
	router.HandleFunc("/execute", Chain(ExecuteHandler, Method("POST"), Logging()))
	router.HandleFunc("/ws", Chain(WebsocketHandler, Logging()))

	subrouter := router.PathPrefix("/model").Subrouter()
	subrouter.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(CurrentModelsHandler)))

	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir(static))))
	return router
}

// Start new server
func Start() *http.Server {

	router := InitRouter()

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return srv
}

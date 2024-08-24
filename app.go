package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type App struct {
	Router *mux.Router
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Print(fmt.Sprintf("Server running on port [%s]", addr))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Connection", "Priority", "Sec-GPC", "DNT", "Referer", "User-Agent", "Host", "Accept", "Accept-Language", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "Origin", "Cache-Control", "X-Requested-With"},
		AllowedMethods:   []string{"GET", "OPTIONS", "POST", "PUT", "DELETE"},
	})

	hnd := c.Handler(a.Router)
	log.Fatal(http.ListenAndServe(addr, hnd))
}

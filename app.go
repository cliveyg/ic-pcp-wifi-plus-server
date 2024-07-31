package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
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
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

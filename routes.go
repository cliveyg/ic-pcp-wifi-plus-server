package main

import log "github.com/sirupsen/logrus"

func (a *App) initializeRoutes() {

	// endpoints
	log.Debug("-----------------------------")
	a.Router.HandleFunc("/test", a.testTings).Methods("GET")
	a.Router.HandleFunc("/wpa/status", a.getWPACliStatus).Methods("GET")
	a.Router.HandleFunc("/system/{action}", a.systemAction).Methods("GET")
	a.Router.HandleFunc("/wifi/{action}", a.wifiAction).Methods("GET")
}

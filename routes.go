package main

import log "github.com/sirupsen/logrus"

func (a *App) initializeRoutes() {

	// endpoints
	log.Debug("In initializeRoutes")
	a.Router.HandleFunc("/test", a.testTings).Methods("GET")
	a.Router.HandleFunc("/picore", a.getPiCoreDetails).Methods("GET")
	a.Router.HandleFunc("/status", a.getSystemStatus).Methods("GET")
	a.Router.HandleFunc("/wifi-status", a.getWifiStatus).Methods("GET")
	a.Router.HandleFunc("/wifi-ssid", a.getWifiSSID).Methods("GET")
}

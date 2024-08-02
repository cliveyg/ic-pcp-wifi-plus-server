package main

func (a *App) initializeRoutes() {

	// endpoints
	a.Router.HandleFunc("/status", a.getSystemStatus).Methods("GET")
	a.Router.HandleFunc("/wifi-status", a.getWifiStatus).Methods("GET")
	a.Router.HandleFunc("/wifi-ssid", a.getWifiSSID).Methods("GET")
}

package main

func (a *App) initializeRoutes() {

	// endpoints
	a.Router.HandleFunc("/test", a.testTings).Methods("GET")
	a.Router.HandleFunc("/wpa/status", a.getWPACliStatus).Methods("GET")
	a.Router.HandleFunc("/system/picore", a.getPiCoreDetails).Methods("GET")
	a.Router.HandleFunc("/system/reboot", a.RebootSystem).Methods("GET")
	a.Router.HandleFunc("/system/status", a.getSystemStatus).Methods("GET")
	//a.Router.HandleFunc("/wifi/status", a.getWifiStatus).Methods("GET")
	//a.Router.HandleFunc("/wifi/ssid", a.getWifiSSID).Methods("GET")
	a.Router.HandleFunc("/wifi/{action}", a.wifiAction).Methods("GET")
}

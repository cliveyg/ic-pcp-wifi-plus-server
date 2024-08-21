package main

func (a *App) initializeRoutes() {

	// endpoints
	a.Router.HandleFunc("/test", a.testTings).Methods("GET", "POST", "OPTIONS")
	a.Router.HandleFunc("/wpa/status", a.getWPACliStatus).Methods("GET", "OPTIONS")
	a.Router.HandleFunc("/wifiplus/switcher", a.wpSwitcher).Methods("GET", "OPTIONS")
	a.Router.HandleFunc("/system/{action}", a.systemAction).Methods("GET", "PUT", "OPTIONS")
	a.Router.HandleFunc("/wifi", a.wifiSwitchNetwork).Methods("POST", "DELETE", "OPTIONS")
	a.Router.HandleFunc("/wifi/{action}", a.wifiAction).Methods("GET", "OPTIONS")
	a.Router.HandleFunc("/wap/{action}", a.wapAction).Methods("GET", "PUT", "OPTIONS")
	a.Router.HandleFunc("/wap", a.wapInfo).Methods("GET", "OPTIONS")
	a.Router.HandleFunc("/wap", a.wapAddRemove).Methods("POST", "DELETE", "OPTIONS")
	a.Router.HandleFunc("{rest:[a-zA-Z0-9=\\-\\/]+}", a.return404).Methods("GET", "OPTIONS")
}

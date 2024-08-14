package main

func (a *App) initializeRoutes() {

	// endpoints
	a.Router.HandleFunc("/test", a.testTings).Methods("GET")
	a.Router.HandleFunc("/wpa/status", a.getWPACliStatus).Methods("GET")
	a.Router.HandleFunc("/wifiplus/switcher", a.wpSwitcher).Methods("GET")
	a.Router.HandleFunc("/system/{action}", a.systemAction).Methods("GET", "PUT")
	a.Router.HandleFunc("/wifi", a.wifiSwitchNetwork).Methods("POST", "DELETE")
	a.Router.HandleFunc("/wifi/{action}", a.wifiAction).Methods("GET")
	a.Router.HandleFunc("/wap/{action}", a.wapAction).Methods("GET", "PUT")
	a.Router.HandleFunc("/wap", a.wapInfo).Methods("GET")
	a.Router.HandleFunc("/wap", a.wapAddRemove).Methods("POST", "DELETE")
	a.Router.HandleFunc("{rest:[a-zA-Z0-9=\\-\\/]+}", a.return404).Methods("GET")
}

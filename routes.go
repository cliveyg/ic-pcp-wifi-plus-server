package main

func (a *App) initializeRoutes() {

	// endpoints
	a.Router.HandleFunc("/test", a.testTings).Methods("GET")
	a.Router.HandleFunc("/wpa/status", a.getWPACliStatus).Methods("GET")
	a.Router.HandleFunc("/system/{action}", a.systemAction).Methods("GET")
	a.Router.HandleFunc("/wifi/{action}", a.wifiAction).Methods("GET")
	a.Router.HandleFunc("/wap/{action}", a.wifiAction).Methods("GET")
	a.Router.HandleFunc("/wap/{action}", a.wapAction).Methods("GET", "POST", "PUT", "DELETE")
	a.Router.HandleFunc("{rest:[a-zA-Z0-9=\\-\\/]+}", a.return404).Methods("GET")
}

package main

func (a *App) initializeRoutes() {

	// endpoints
	a.Router.HandleFunc("/status", a.getStatus).Methods("GET")

}

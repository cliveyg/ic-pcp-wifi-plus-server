package main

import (
	"io"
	"log"
	"net/http"
)

// ----------------------------------------------------------------------------

func (a *App) getStatus(w http.ResponseWriter, _ *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	mess := `{"message": "System running..."}`
	if _, err := io.WriteString(w, mess); err != nil {
		log.Fatal(err)
	}

}

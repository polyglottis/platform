// Package server contains all server features which are common to frontend and api.
package server

import (
	"log"
	"net/http"
)

func internalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error. Please try again in a few seconds.\n"))
}

func InternalError(r *http.Request, w http.ResponseWriter, err error) {
	log.Printf("Internal server error: %v (while serving %s)", err, r)
	internalError(w)
}

func Recovered(r *http.Request, w http.ResponseWriter, rec interface{}) {
	log.Printf("Panic! Recovered from %v while serving %v", rec, r)
	internalError(w)
}

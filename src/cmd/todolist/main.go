package main

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func Health(writer http.ResponseWriter, request *http.Request) {
	log.Info("Server is UP")

	writer.Header().Set("Content-Type", "application/json")
	io.WriteString(writer, `{"health": UP}`)
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	log.Info("Starting server")

	router := mux.NewRouter()
	router.HandleFunc("/health", Health).Methods("GET")
	http.ListenAndServe(":8000", router)
}

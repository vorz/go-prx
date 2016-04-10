package main

import (
	"net/http"
)

func indexHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("PRIVET"))
}

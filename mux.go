package main

import "net/http"

func NewMux() http.Handler {
	mux := http.NewServerMux()
	mux.HandlerFunc()
}

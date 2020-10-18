package main

import (
	"log"
	"net/http"
	"os"
)

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://127.0.0.1:"+os.Getenv("PORT")+r.RequestURI, http.StatusMovedPermanently)
}

func redirectServer() {
	if err := http.ListenAndServe(":80", http.HandlerFunc(redirectTLS)); err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
func init() {
	go redirectServer()
}

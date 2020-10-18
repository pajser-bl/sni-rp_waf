package main

import (
	"io"
	"log"
	"os"
	"strings"
)

// Log the environment variables for reverse proxy
func logSetup() {
	log.Printf("======================================================\n")
	log.Printf("Proxy server will run on: %s\n", os.Getenv("PORT"))
	log.Printf("======================================================\n")
	for _, t := range strings.Split(os.Getenv("TARGETS"), " ") {
		log.Printf("Redirecting: %s to http://localhost:%s\n", strings.Split(t, ":")[0], strings.Split(t, ":")[1])
	}
	log.Printf("======================================================\n")
	log.Printf("Server will handle " + os.Getenv("REQUESTS_RATE") + " requests per second.")
	log.Printf("Server will handle " + os.Getenv("REQUESTS_RATE_PER_IP") + " requests per IP per second.")
	log.Printf("Max request header size is %s Bytes.", os.Getenv("MAX_HEADER_BYTES"))
	log.Printf("Read header request timeout is %s seconds.", os.Getenv("READ_HEADER_TIMEOUT"))
	log.Printf("Read request timeout is %s seconds.", os.Getenv("READ_TIMEOUT"))
	log.Printf("Write request timeout is %s seconds.", os.Getenv("WRITE_TIMEOUT"))
	log.Printf("Idle request timeout is %s seconds.", os.Getenv("IDLE_TIMEOUT"))
	log.Printf("======================================================\n")
	if os.Getenv("WAP") == "true" {
		log.Printf("Web application firewal is ON.\n")
	} else {
		log.Printf("Web application firewal is OFF.\n")
	}
	log.Printf("======================================================\n")
	if len(w) > 0 {
		log.Printf("Whitelist:")
		for i, r := range w {
			log.Printf("[%s] => %s.", i, r)
		}
		log.Printf("======================================================\n")
	}

}

// Log the typeform payload and redirection url
func logRequestPayload(r string, p string, i string) {
	log.Printf("Redirecting[%s]: %s -> %s\n", i, r, p)
}

// Set up for logs to be printed to stdout and file.
func init() {
	f, err := os.OpenFile("proxy.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	//defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
}

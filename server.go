package main

import (
	"bytes"
	"crypto/tls"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func getProxyUrl(requestUrl string) string {
	targets := strings.Split(os.Getenv("TARGETS"), " ")
	for _, t := range targets {
		r := strings.Split(t, ":")
		name, port := r[0], r[1]
		if strings.HasPrefix(requestUrl, "/"+strings.ToLower(name)) {
			return "http://localhost:" + port + strings.Split(requestUrl, "/"+strings.ToLower(name))[1]
		}
	}
	return "http://localhost:" + os.Getenv("DEFAULT_TARGET")
}

type transport struct {
	http.RoundTripper
}

var _ http.RoundTripper = &transport{}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	log.Println(resp)
	return resp, nil
}

// Serve a reverse proxy for given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	tUrl, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(tUrl)
	//proxy.Transport = &transport{http.DefaultTransport}
	//_, _ = proxy.Transport.RoundTrip(req)
	proxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Header.Get("Host"))
		req.Header.Add("X-Origin-Host", req.Host)
		req.URL.Scheme = tUrl.Scheme
		req.URL.Host = tUrl.Host
		req.URL.Path = tUrl.Path
		req.RequestURI = tUrl.Path
	}
	proxy.ServeHTTP(res, req)

}

// On request, forward it to the appropriate url
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	reqUrl := req.URL.String()
	u := getProxyUrl(reqUrl)

	b, _ := ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewReader(b))

	if os.Getenv("WAP") == "true" {
		if isSafe, msg := isSafe(req); !isSafe {
			log.Printf("Request from %s discarded. [%s]\n", req.RemoteAddr, msg)
			res.WriteHeader(http.StatusForbidden)
			return
		} else {
			log.Printf("Request from %s is [%s].\n", req.RemoteAddr, msg)
			req.Body = ioutil.NopCloser(bytes.NewReader(b))
		}
	}
	logRequestPayload(reqUrl, u, req.RemoteAddr)
	serveReverseProxy(u, res, req)
}

func getServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRequestAndRedirect)
	port := ":" + getEnv("PORT", "443")
	maxHeaderBytes, _ := strconv.Atoi(getEnv("MAX_HEADER_BYTES", "8190"))
	readHeaderTimeout, _ := strconv.Atoi(getEnv("READ_HEADER_TIMEOUT", "1"))
	readTimeout, _ := strconv.Atoi(getEnv("READ_TIMEOUT", "5"))
	writeTimeout, _ := strconv.Atoi(getEnv("WRITE_TIMEOUT", "10"))
	idleTimeout, _ := strconv.Atoi(getEnv("IDLE_TIMEOUT", "120"))
	certFile := getEnv("CERTIFICATE", "server.crt")
	keyFile := getEnv("CERTIFICATE_KEY", "server.key")
	cert, _ := tls.LoadX509KeyPair(certFile, keyFile)
	return &http.Server{
		Addr:              port,
		Handler:           limit(mux),
		MaxHeaderBytes:    maxHeaderBytes,
		ReadHeaderTimeout: time.Duration(readHeaderTimeout) * time.Second,
		ReadTimeout:       time.Duration(readTimeout) * time.Second,
		WriteTimeout:      time.Duration(writeTimeout) * time.Second,
		IdleTimeout:       time.Duration(idleTimeout) * time.Second,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				cert,
			},
		},
	}
}

// Runs before main, load environment variables from .env
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found!")
		log.Print("Please make one containing PORT, CERTIFICATE, CERTIFICATE_KEY, TARGETS, and set wap ENABLED...")
	}
}

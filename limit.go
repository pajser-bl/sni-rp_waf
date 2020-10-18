package main

import (
	"golang.org/x/time/rate"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var rR, _ = strconv.Atoi(os.Getenv("REQUESTS_RATE"))
var rIP, _ = strconv.Atoi(os.Getenv("REQUESTS_RATE_PER_IP"))

var limiter = rate.NewLimiter(1000, 1000*5)

func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Global limit
		if limiter.Allow() == false {
			http.Error(writer, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		// Concurrent limit
		c, err := strconv.Atoi(os.Getenv("CONCURRENCY_LIMIT"))
		if err != nil {
			c = 100
		}
		if len(visitors) > c {
			http.Error(writer, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		// Per user(IP) limit
		ip, _, err := net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			log.Println(err.Error())
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		limiter := getVisitor(ip)
		if limiter.Allow() == false {
			http.Error(writer, http.StatusText(429), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	vis, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(100, 100*5)
		visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}
	vis.lastSeen = time.Now()
	return vis.limiter
}

// Cleans all visitors that have not been seen in last 3 minutes.
func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

// Will run cleanup goroutine
func init() {
	go cleanupVisitors()
}

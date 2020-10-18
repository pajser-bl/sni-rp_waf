package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type bE struct {
	t time.Time
	r string
}

var w = make(map[string]string)
var b = make(map[string]*bE)
var m sync.Mutex
var bCT = 4320

func isSafe(r *http.Request) (bool, string) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if e, ok := w[ip]; ok {
		return true, "Whitelisted(ip: " + ip + ", reason:" + e + ")"
	}
	if e, ok := b[ip]; ok {
		return false, "Blacklisted(" + e.r + "@" + e.t.Format("2006-01-02 15:04:05") + ")"
	}
	if v, ok := os.LookupEnv("BLACKLIST_TIME"); ok {
		bCT, _ = strconv.Atoi(v)
	}
	d, _ := ioutil.ReadAll(r.Body)
	s := string(d)
	q := r.URL.RawQuery
	if isSQLInjection(s) || isSQLInjection(q) {
		b[ip] = &bE{time.Now(), "SQL injection"}
		return false, "SQL injection"
	}
	if isXSSInjection(s) || isXSSInjection(q) {
		b[ip] = &bE{time.Now(), "XSS injection"}
		return false, "XSS injection"
	}
	co, _ := r.Cookie("sni")
	if isCookiePoisoned(ip, co) {
		b[ip] = &bE{time.Now(), "Cookie poisoning"}
		return false, "Cookie poisoning"
	}
	return true, "OK"
}
func init() {
	for _, v := range strings.Split(os.Getenv("WHITELIST"), " ") {
		if v != "" {
			i := strings.Split(v, ":")[0]
			r := strings.Split(v, ":")[1]
			w[i] = r
		}
	}
	go cleanBlacklist()
}

func cleanBlacklist() {
	for {
		time.Sleep(time.Minute)
		m.Lock()
		for i, v := range b {
			if time.Since(v.t) > time.Minute*time.Duration(bCT) {
				delete(b, i)
			}
		}
		m.Unlock()
	}
}

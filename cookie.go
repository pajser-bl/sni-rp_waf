package main

import (
	"log"
	"net/http"
)

var c = make(map[string]string)

func isCookiePoisoned(ip string, co *http.Cookie) bool {
	if e, ok := c[ip]; ok {
		log.Println(e, co.Value)
		if e == co.Value {
			return true
		}
	}
	return false
}

func saveCookie(ip string, co *http.Cookie) {
	if co != nil {
		c[ip] = co.Value
		log.Println(ip, co.Value)
	}
}

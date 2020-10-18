package main

import (
	"os"
	"regexp"
	"strings"
)

var x []regexp.Regexp

func init() {
	for _, v := range strings.Split(os.Getenv("XSS_REGEX"), " ") {
		if re, err := regexp.Compile(v); err == nil {
			x = append(x, *re)
		}
	}
}

func isXSSInjection(s string) bool {
	for _, v := range x {
		if v.MatchString(s) {
			return true
		}
	}
	return false
}

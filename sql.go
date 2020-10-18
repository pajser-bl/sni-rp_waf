package main

import (
	"os"
	"regexp"
	"strings"
)

var r []regexp.Regexp

func init() {
	for _, v := range strings.Split(os.Getenv("SQL_REGEX"), " ") {
		if re, err := regexp.Compile(v); err == nil {
			r = append(r, *re)
		}
	}
}

func isSQLInjection(s string) bool {
	for _, v := range r {
		if v.MatchString(s) {
			return true
		}
	}
	return false
}

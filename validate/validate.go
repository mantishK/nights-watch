package validate

import (
	"regexp"
	"strings"

	"net/url"
)

func Password(password string) bool {
	if len(password) < 6 {
		return false
	}
	return true
}

func Email(email string) bool {
	if strings.Contains(email, " ") {
		return false
	}
	exp, err := regexp.Compile(`[a-zA-Z0-9\+\.\_\%\-\+]{1,256}\@[a-zA-Z0-9][a-zA-Z0-9\-]{0,64}(\.[a-zA-Z0-9][a-zA-Z0-9\-]{0,25})+`)
	if err != nil {
		return false
	}
	if !exp.MatchString(email) {
		return false
	}
	return true
}

func Url(u string) bool {
	if _, err := url.Parse(u); err != nil {
		return false
	}
	return true
}

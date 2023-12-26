package common

import (
	"errors"
	"net/url"
)

func ValidateInput(urlString, method, data *string) error {
	if *urlString == "" {
		return errors.New("a URL must be provided with the -u flag")
	}
	if _, err := url.ParseRequestURI(*urlString); err != nil {
		return errors.New("invalid URL provided")
	}
	if *method != "GET" && *method != "POST" {
		return errors.New("HTTP Method must be GET or POST")
	}
	if *method == "POST" && *data == "" {
		return errors.New("data must be provided for POST requests with the -d flag")
	}
	return nil
}

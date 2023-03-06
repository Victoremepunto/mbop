package handlers

import (
	"fmt"
	"regexp"
)

const CERT_HEADER = "x-rh-certauth-cn"

var cnMatcher = regexp.MustCompile(`/CN=(.*)$`)

func getCertCN(header string) (string, error) {
	if header == "" || !cnMatcher.MatchString(header) {
		return "", fmt.Errorf("[x-rh-certauth-cn] header not present")
	}

	parts := cnMatcher.FindAllStringSubmatch(header, 1)
	return parts[0][1], nil
}

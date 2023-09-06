package url

import (
	"errors"
	"fmt"
	"strings"
)

type Url struct {
	Scheme   string
	Hostname string
	Path     string
	port     string
}

func (u *Url) GetPort() string {
	return u.port
}

func (u *Url) String() string {
	return fmt.Sprintf("Host: %s Scheme: %s Port: %s", u.Hostname, u.Scheme, u.port)
}

func schemeIsInvalid(data string) bool {
	return strings.Index(data, "://") < 0
}

func hostnameIsInvalid(data string) bool {
	if data == "" || strings.Index(data, ".") < 0 {
		return true
	} else {
		return false
	}
}

func Parse(data string) (*Url, error) {
	scheme, rest, ok := parseScheme(data)
	if !ok {
		return nil, errors.New("missing scheme")
	}
	hostname, path, ok := parseHostnamePath(rest)
	if !ok {
		return nil, errors.New("invalid hostname")
	}
	portData := parsePort(hostname)
	parsed := Url{
		Scheme:   scheme,
		Hostname: hostname,
		Path:     path,
		port:     portData,
	}
	return &parsed, nil
}

func parseHostnamePath(data string) (hostname, path string, ok bool) {
	var pathData string
	hostData := strings.Split(data, "/")
	if len(hostData) > 1 {
		pathData = strings.Join(hostData[1:], "/")
	}
	if hostnameIsInvalid(hostData[0]) {
		return "", "", false
	}
	return hostData[0], pathData, true
}

func parseScheme(data string) (scheme, rest string, ok bool) {
	if schemeIsInvalid(data) {
		return "", "", false
	}
	schemeRest := strings.Split(data, "://")
	return schemeRest[0], schemeRest[1], true
}

func parsePort(hostname string) (port string) {
	if strings.Index(hostname, ":") > 1 {
		return strings.Split(hostname, ":")[1]
	} else {
		return ""
	}

}

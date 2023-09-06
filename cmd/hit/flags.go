package main

import (
	"errors"
	"example.com/cmll/url"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const usageText = `
Usage:
	hit [options] url
Options:
`

type flags struct {
	url     string
	n, c    number
	timeout time.Duration
	method  option
	headers headers
}

type number int

func (n *number) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	switch {
	case err != nil:
		err = errors.New("parse error")
	case v <= 0:
		err = errors.New("must be positive")
	}
	*n = number(v)
	return err
}

func (n *number) String() string {
	return strconv.Itoa(int(*n))
}

type option string

func (o *option) Set(s string) error {
	validOptions := []string{"GET", "POST", "PUT"}
	valid := false
	for _, option := range validOptions {
		if strings.Contains(s, option) {
			valid = true
		}
	}
	if s == "" {
		valid = true
	}
	if !valid {
		return errors.New("invalid method " + s)
	} else {
		*o = option(s)
		return nil
	}
}

func (o *option) String() string {
	return string(*o)
}

type headers []string

func (h *headers) Set(s string) error {
	*h = append(*h, s)
	return nil
}

func (h *headers) String() string {
	return strings.Join(*h, " ")
}

func (f *flags) parse(s *flag.FlagSet, args []string) error {
	s.Usage = func() {
		fmt.Fprintln(os.Stderr, usageText[1:])
		flag.PrintDefaults()
	}
	s.Var(&f.n, "n", "Number of requests to make")
	s.Var(&f.c, "c", "Concurrency level")
	s.DurationVar(&f.timeout, "t", time.Duration(0), "Timeout value")
	s.Var(&f.method, "m", "Method, must be one of GET, POST, PUT")
	s.Var(&f.headers, "H", "Headers for the request, can be multiple")
	if err := s.Parse(args); err != nil {
		fmt.Println(s.Output(), err)
		return err
	}
	// First positional argument is url
	f.url = s.Arg(0)
	if err := f.validate(); err != nil {
		fmt.Fprintln(s.Output(), err)
		return err
	}
	return nil
}

func (f *flags) validate() error {
	err := validateUrl(f.url)
	if err != nil {
		return err
	}
	if f.c > f.n {
		return fmt.Errorf("-c=%d should be less or equal than -n=%d", f.c, f.n)
	}
	return nil
}

func validateUrl(s string) error {
	_, err := url.Parse(s)
	switch {
	case strings.TrimSpace(s) == "":
		{
			err = errors.New("argument url is required")
		}
	}
	if err != nil {
		return err
	}
	return nil
}

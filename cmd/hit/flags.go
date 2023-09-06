package main

import (
	"errors"
	"example.com/cmll/url"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const usageText = `
Usage:
	hit [options] url
Options:
`

type flags struct {
	url  string
	n, c int
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

func toNumber(p *int) *number {
	return (*number)(p)
}

func (f *flags) parse(s *flag.FlagSet, args []string) error {
	s.Usage = func() {
		fmt.Fprintln(os.Stderr, usageText[1:])
		flag.PrintDefaults()
	}
	s.Var(toNumber(&f.n), "n", "Number of requests to make")
	s.Var(toNumber(&f.c), "c", "Concurrency level")
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

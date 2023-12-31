package main

import (
	"bytes"
	"flag"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

type testEnv struct {
	args           string
	stdout, stderr bytes.Buffer
}

func (e *testEnv) run() error {
	s := flag.NewFlagSet("hit", flag.ContinueOnError)
	s.SetOutput(&e.stderr)
	return run(s, strings.Fields(e.args), &e.stdout)
}

func TestRun(t *testing.T) {
	t.Parallel()
	numCpu := strconv.Itoa(runtime.NumCPU())

	happy := map[string]struct{ in, out string }{
		"url": {
			"https://google.com",
			fmt.Sprintf("Making %s requests to https://google.com with a concurrency of %s", numCpu, numCpu),
		},
		"n_c": {
			"-n=10 -c=5 https://google.com",
			"Making 10 requests to https://google.com with a concurrency of 5",
		},
		"timeout": {
			"-t=5s https://foo.com",
			"Making 12 requests to https://foo.com with a concurrency of 12 and timeout of 5s",
		},
		"method": {
			"-m GET https://foo.com",
			"Making 12 GET requests to https://foo.com with a concurrency of 12",
		},
		"headers": {
			"-H=Accept:text/json https://foo.com",
			"Headers: Accept:text/json",
		},
		"headers/mult": {
			"-H=Accept:text/json -H=User-Agent:hit https://foo.com",
			"Headers: Accept:text/json User-Agent:hit",
		},
	}
	for name, tt := range happy {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := &testEnv{args: tt.in}
			if err := e.run(); err != nil {
				t.Errorf("\rgot:\n %q\nwant:\n nil err", err)
			}
			if out := e.stdout.String(); !strings.Contains(out, tt.out) {
				t.Errorf("\ngot:\n %s\nwant:\n %s", out, tt.out)
			}
		})
	}
}

func TestRunBad(t *testing.T) {
	bad := map[string]string{
		"url/missing":        "",
		"url/bad":            "https://something_bad",
		"url/missing_schema": "something_bad",
		"c/negative":         "-c -1 https://foo.com",
		"c/err":              "-c bad https://foo.com",
		"n/negative":         "-n -1 https://google.com",
		"n/err":              "-n bad https://foo.com",
		"n/c":                "-c 10 -n 2 https://foo.com",
		"c/zero":             "-c 0 https://foo.com",
		"n/zero":             "-n 0 https://foo.com",
		"t/negative":         "-t -1 https://foo.com",
		"t/bad":              "-t bad https://foo.com",
		"m/bad":              "-m OPTIONS https://foo.com",
	}
	for name, tt := range bad {
		tt := tt
		t.Run(name, func(t *testing.T) {

			e := &testEnv{args: tt}
			if err := e.run(); err == nil {
				t.Error("got nil; want err")
			}
			if e.stderr.Len() == 0 {
				t.Error("stderr = 0 bytes; want > 0")
			}
		})
	}
}

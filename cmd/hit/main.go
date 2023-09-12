package main

import (
	"flag"
	"fmt"
	"github.com/cmll/hit/hit"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"
)

const (
	bannerText = `
 __  __     __     ______
/\ \_\ \   /\ \   /\__  _\
\ \  __ \  \ \ \  \/_/\ \/
 \ \_\ \_\  \ \_\    \ \_\
  \/_/\/_/   \/_/     \/_/
`
)

func main() {
	fmt.Fprint(os.Stdout, bannerText)
	if err := run(flag.CommandLine, os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error occurred:", err)
		os.Exit(1)
	}
}

func run(s *flag.FlagSet, args []string, out io.Writer) error {
	f := &flags{
		n: number(runtime.NumCPU()),
		c: number(runtime.NumCPU()),
	}
	if err := f.parse(s, args); err != nil {
		return err
	}

	var timeoutText string
	if f.timeout != time.Duration(0) {
		timeoutText = fmt.Sprintf(" and timeout of %v", f.timeout)
	}

	var methodText string
	if f.method != "" {
		methodText = fmt.Sprintf(" %s", f.method)
	}

	if f.headers != nil {
		fmt.Fprintf(out, "Headers: %s\n", f.headers.String())
	}
	fmt.Fprintf(out, "Making %d%s requests to %s with a concurrency of %d%s", f.n, methodText, f.url, f.c, timeoutText)

	request, err := http.NewRequest(http.MethodGet, f.url, http.NoBody)
	if err != nil {
		return err
	}
	var c hit.Client
	sum := c.Do(request, int(f.n))
	sum.Finalize(2 * time.Second)
	sum.FPrint(out)

	return nil
}

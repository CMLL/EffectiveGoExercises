package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/cmll/hit/hit"
	"io"
	"net/http"
	"os"
	"os/signal"
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

	const timeout = time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	// Cancel and stop functions are done to tell a context is done.
	// We defer on this so that when the program gracefully ends, the contexts stop.
	defer stop()
	defer cancel()

	request, err := http.NewRequest(http.MethodGet, f.url, http.NoBody)
	if err != nil {
		return err
	}
	c := &hit.Client{
		C:       int(f.c),
		RPS:     f.rps,
		Timeout: 10 * time.Second,
	}
	sum := c.Do(ctx, request, int(f.n))
	sum.FPrint(out)

	if err := ctx.Err(); errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("timed out in %v", timeout)
	}

	return nil
}

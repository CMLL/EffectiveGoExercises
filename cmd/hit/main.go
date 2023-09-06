package main

import (
	"flag"
	"fmt"
	"io"
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
		os.Exit(1)
	}
}

func run(s *flag.FlagSet, args []string, out io.Writer) error {
	f := &flags{
		n: runtime.NumCPU(),
		c: runtime.NumCPU(),
	}
	if err := f.parse(s, args); err != nil {
		return err
	}

	var timeoutText string
	if f.timeout != time.Duration(0) {
		timeoutText = fmt.Sprintf(" and timeout of %v", f.timeout)
	}

	fmt.Fprintf(out, "Making %d requests to %s with a concurrency of %d%s", f.n, f.url, f.c, timeoutText)

	return nil
}

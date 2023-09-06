package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
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

	fmt.Fprintf(out, "Making %d requests to %s with a concurrency of %d", f.n, f.url, f.c)

	return nil
}

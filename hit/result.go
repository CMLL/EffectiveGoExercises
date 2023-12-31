package hit

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type Result struct {
	RPS      float64
	Requests int
	Errors   int
	Bytes    int64
	Duration time.Duration
	Fastest  time.Duration
	Slowest  time.Duration
	Status   int
	Success  float64
	Error    error
}

func (r *Result) Merge(data *Result) {
	r.Requests++
	r.Bytes += data.Bytes

	if r.Fastest == 0 || data.Duration < r.Fastest {
		r.Fastest = data.Duration
	}
	if data.Duration > r.Slowest {
		r.Slowest = data.Duration
	}

	switch {
	case data.Error != nil:
		fallthrough
	case data.Status >= 400:
		r.Errors++
	}
}

func (r *Result) Finalize(total time.Duration) *Result {
	r.Duration = total
	r.RPS = float64(r.Requests) / r.Duration.Seconds()
	r.Success = r.success()
	return r
}

func (r *Result) success() float64 {
	successReq := r.Requests - r.Errors
	return float64(successReq) * float64(100) / float64(r.Requests)
}

func (r *Result) String() string {
	var s strings.Builder
	r.FPrint(&s)
	return s.String()
}

func (r *Result) FPrint(out io.Writer) {
	fmt.Fprintf(out, `
Summary:
	Requests: %d
	Errors: %d
	Success: %.0f%%
	RPS: %.1f
	Duration: %v
	Fastest: %v
	Slowest: %v
`, r.Requests, r.Errors, r.success(), r.RPS, r.Duration, r.Fastest, r.Slowest)
}

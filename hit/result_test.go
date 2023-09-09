package hit

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestResult_MergeMultipleRequests(t *testing.T) {
	sut := Result{}

	newResult := Result{
		Bytes: 256,
	}
	sut.Merge(&newResult)

	assert.Equal(t, 1, sut.Requests)
	assert.Equal(t, int64(256), sut.Bytes)

	newResult = Result{
		Bytes: 256,
	}
	sut.Merge(&newResult)

	assert.Equal(t, 2, sut.Requests)
	assert.Equal(t, int64(512), sut.Bytes)
}

func TestResult_TracksDuration(t *testing.T) {
	sut := Result{}

	results := []*Result{
		{Duration: time.Duration(30)},
		{Duration: time.Duration(10)},
		{Duration: time.Duration(25)},
	}
	for _, result := range results {
		sut.Merge(result)
	}

	assert.Equal(t, time.Duration(30), sut.Slowest)
	assert.Equal(t, time.Duration(10), sut.Fastest)
}

func TestErrorResultGetsCounted(t *testing.T) {

	testCases := []struct {
		name           string
		results        []*Result
		expectedErrors int
	}{
		{"50%", []*Result{
			{Error: errors.New("bad result")},
			{Status: 200},
		}, 1},
	}
	for _, tc := range testCases {
		sut := Result{}
		t.Run(tc.name, func(t *testing.T) {
			for _, result := range tc.results {
				sut.Merge(result)
				assert.Equal(t, tc.expectedErrors, sut.Errors)
			}
		})
	}
}

func TestFinalizeSetsTotalDuration(t *testing.T) {
	sut := Result{}

	result := sut.Finalize(time.Duration(10))

	assert.Equal(t, time.Duration(10), result.Duration, "want 10s got %v", result.Duration)
}

func TestFinalizeCalculatesRPS(t *testing.T) {
	sut := Result{}

	testCases := []struct {
		name     string
		results  []*Result
		expected float64
	}{
		{"1x1", []*Result{{Duration: time.Second}}, 1},
		{"3x1", []*Result{{Duration: 3 * time.Second}, {Duration: 3 * time.Second}, {Duration: 3 * time.Second}}, 0.4444444444444444},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var totalDuration time.Duration
			for _, result := range tc.results {
				totalDuration += result.Duration
				sut.Merge(result)
			}
			sut.Finalize(totalDuration)
			assert.Equal(t, tc.expected, sut.RPS)
		})
	}
}

func TestFinalizeCalculatesSuccessRate(t *testing.T) {
	testCases := []struct {
		name     string
		requests int
		errors   int
		expected float64
	}{
		{"50%", 10, 5, 50},
		{"25%", 10, 7, 30},
		{"90%", 10, 1, 90},
	}
	for _, tc := range testCases {
		sut := Result{
			Requests: tc.requests, Errors: tc.errors,
		}
		sut.Finalize(10)
		assert.Equal(t, tc.expected, sut.Success)
	}
}

func TestResultPrintsWriter(t *testing.T) {
	sut := Result{
		Requests: 4,
		Errors:   2,
	}
	buffer := bytes.Buffer{}

	sut.Fprint(&buffer)

	assert.NotEmpty(t, buffer)
	data := buffer.String()
	assert.Contains(t, data, "Summary")
	assert.Contains(t, data, "Requests: 4")
	assert.Contains(t, data, "Errors: 2")
	assert.Contains(t, data, "Success: 50.0")
}

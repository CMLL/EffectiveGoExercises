package url_test

import (
	"github.com/cmll/hit/url"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	subject = "https://google.com"
)

func TestParseReturnsNoError(t *testing.T) {

	_, err := url.Parse(subject)
	assert.NoError(t, err, "Expected nil but got %q", err)
}

func TestParseErrorCases(t *testing.T) {
	testCases := []struct {
		target   string
		expected string
	}{
		{
			"google.com",
			"missing scheme",
		},
		{
			"https:///bad",
			"bad url",
		},
		{
			"https://somethingotherwithoutdot",
			"bad url",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.target, func(t *testing.T) {
			sut, err := url.Parse(tc.target)
			assert.Nil(t, sut)
			assert.Errorf(t, err, tc.expected)
		})
	}
}

func TestParseReturnsUrlScheme(t *testing.T) {
	testCases := []struct {
		name     string
		target   string
		expected string
	}{
		{
			"https",
			"https://google.com",
			"https",
		},
		{
			"http",
			"http://gogle.com",
			"http",
		},
		{
			"ftp",
			"ftp://google.com",
			"ftp",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut, _ := url.Parse(tc.target)
			assert.Equal(t, sut.Scheme, tc.expected, "%s Expected %v got %v", tc.name, tc.expected, sut.Scheme)
		})
	}
}

func TestParserReturnsHostname(t *testing.T) {
	testCases := []struct {
		target   string
		expected string
		port     string
	}{
		{
			"https://google.com",
			"google.com",
			"",
		},
		{
			"https://facebook.com",
			"facebook.com",
			"",
		},
		{
			"https://reddit.com/path",
			"reddit.com",
			"",
		},
		{
			"https://google.com:443",
			"google.com:443",
			"443",
		},
		{
			"https://127.0.0.1:443/something",
			"127.0.0.1:443",
			"443",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.target, func(t *testing.T) {
			sut, _ := url.Parse(tc.target)
			assert.Equal(t, tc.expected, sut.Hostname, "Expected %v got %v", tc.expected, sut.Hostname)
			assert.Equal(t, tc.port, sut.GetPort(), "Expected port %v got %v", sut.GetPort(), tc.port)
		})
	}
}

func TestParserReturnsPath(t *testing.T) {
	testCases := []struct {
		target   string
		expected string
	}{
		{
			"https://reddit.com/myPath",
			"myPath",
		},
		{
			"https://google.com/",
			"",
		},
		{
			"ftp://reddit.com/newPath/extra",
			"newPath/extra",
		},
		{
			"ftp://10.0.0.1/ipPath",
			"ipPath",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.target, func(t *testing.T) {
			sut, _ := url.Parse(tc.target)
			assert.Equal(t, sut.Path, tc.expected, "Expected %v got %v", tc.expected, sut.Path)
		})
	}
}

func BenchmarkParse(b *testing.B) {
	b.Logf("Loop %d times", b.N)
	data := "https://reddit.com/myPath"
	for i := 0; i < b.N; i++ {
		url.Parse(data)
	}
	b.ReportAllocs()
}

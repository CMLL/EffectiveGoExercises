package url_test

import (
	"fmt"
	"github.com/cmll/hit/url"
	"log"
)

func ExampleParse() {
	u, err := url.Parse("https://google.com:443/myPath")
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "http"
	fmt.Println(u)
	//	Output:
	//	Host: google.com:443 Scheme: http Port: 443
}

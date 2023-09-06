package url_test

import (
	"example.com/cmll/url"
	"fmt"
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

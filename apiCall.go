package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	url := "http://localhost:8000/api/artifacts"
	method := "GET"

	// payload := strings.NewReader("{\n    \"kargs\": {\n        \"action\": \"list\",\n        \"args\": {}\n    }\n}")
	payload := strings.NewReader("{\n	\"count\": 0\n		\"next\": null,\n		\"previous\": null,\n		\"results\": []	}")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", "Token 6d4989b53848d037ddc5cfbee51559900d3151ef")
	// req.Header.Add("fluencytoken", "b87f0bed-5e50-41e8-781a-e157c9ad9567")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}

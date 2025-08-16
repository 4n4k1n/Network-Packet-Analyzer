package main

import (
	"io"
	"log"
	"net/http"
)

func main() {

	resp, err := http.Get("http://42logtime.com")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	log.Printf(sb)
}

package http

import (
	"io/ioutil"
	"log"
	"net/http"
)

func GetAndReturnBody(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}

func MakeGetRequestAndReturnStatusCode(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}

	return resp.StatusCode, nil
}

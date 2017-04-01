package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func httpGET(url string) []byte {
	resp, err := http.Get(url)
	checkError(err)
	bytes, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	return bytes
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

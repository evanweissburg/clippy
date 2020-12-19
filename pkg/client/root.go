package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func invalid_usage() {
	fmt.Println("Correct usage: <>")
	os.Exit(1)
}

func Execute() {
	if len(os.Args) != 3 {
		invalid_usage()
	}

	switch os.Args[1] {
	case "put":
		filename := os.Args[2]
		fmt.Printf("put %s\n", filename)

		data := "test_data"
		resp, err := http.Post("http://localhost:8080/", "text/plain", bytes.NewBuffer([]byte(data)))
		if err != nil {
			log.Fatalln(err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(string(body))

	case "get":
		clipcode := os.Args[2]
		fmt.Printf("Retrieving clipcode %s\n", clipcode)
		resp, err := http.Get("http://localhost:8080/" + clipcode)
		if err != nil {
			log.Fatalln(err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(string(body))

	default:
		invalid_usage()
	}
}

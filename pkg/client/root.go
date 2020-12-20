package client

import (
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

		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		resp, err := http.Post("http://localhost:8080/", "text/plain", file)
		if err != nil {
			log.Fatalln(err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		clipcode := string(body)

		fmt.Printf("Recieved clipcode %s\n", clipcode)

	case "get":
		clipcode := os.Args[2]
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

package cli

import (
	"fmt"
	"github.com/evanweissburg/clippy/pkg/client"
	"io/ioutil"
	"log"
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

		cl := client.Client{
			BaseURL: "http://localhost:8080/",
		}

		clipcode, err := cl.Upload(file)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("Recieved clipcode %s\n", clipcode)

	case "get":
		clipcode := os.Args[2]

		cl := client.Client{
			BaseURL: "http://localhost:8080/",
		}

		data, err := cl.Download(clipcode)
		if err != nil {
			log.Fatalln(err)
		}

		text, err := ioutil.ReadAll(data)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(string(text))

	default:
		invalid_usage()
	}
}

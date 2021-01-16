package cli

import (
	"fmt"
	"github.com/evanweissburg/clippy/pkg/client"
	"github.com/evanweissburg/clippy/pkg/mnemonic"
	"github.com/mholt/archiver/v3"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
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
		put(filename)

	case "get":
		clipcode := os.Args[2]
		get(clipcode)

	default:
		invalid_usage()
	}
}

func put(filename string) {
	err := archiver.Archive([]string{filename}, ".clip.zip")
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(".clip.zip")
	if err != nil {
		log.Fatal(err)
	}

	clipcode, err := client.Upload("http://localhost:8080/", file)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Recieved clipcode %s\n", clipcode)

	rand.Seed(time.Now().UTC().UnixNano())
	mnemonic, err := mnemonic.CreateSentence(clipcode)
	if err == nil {
		fmt.Printf("Remember it with: %s\n", mnemonic)
	}

	err = os.Remove(".clip.zip")
	if err != nil {
		log.Fatalln(err)
	}
}

func get(clipcode string) {
	data, err := client.Download("http://localhost:8080/", clipcode)
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Create(".clip.zip")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = io.Copy(file, data)
	err = file.Close()
	if err != nil {
		log.Fatalln(err)
	}

	err = archiver.Unarchive(".clip.zip", ".")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.Remove(".clip.zip")
	if err != nil {
		log.Fatalln(err)
	}
}

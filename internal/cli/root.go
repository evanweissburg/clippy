package cli

import (
	"fmt"
	"github.com/evanweissburg/clippy/pkg/client"
	"github.com/evanweissburg/clippy/pkg/mnemonic"
	"github.com/mholt/archiver/v3"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
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
	tempDir, err := ioutil.TempDir("", "clippy-*")
	if err != nil {
		fmt.Printf("Unable to create temporary directory: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	zipFilename := path.Join(tempDir, "clip.zip")

	err = archiver.Archive([]string{filename}, zipFilename)
	if err != nil {
		fmt.Printf("Unable to create zip of %s at %s: %v\n", filename, zipFilename, err)
		return
	}

	zipFile, err := os.Open(zipFilename)
	if err != nil {
		fmt.Printf("Unable to open zip file %s: %v\n", zipFilename, err)
		return
	}
	defer zipFile.Close()

	clipcode, err := client.Upload("http://localhost:8080/", zipFile)
	if err != nil {
		fmt.Printf("Unable to upload to server: %v\n", err)
		return
	}

	fmt.Printf("Recieved clipcode %s\n", clipcode)

	rand.Seed(time.Now().UTC().UnixNano())
	mnemonic, err := mnemonic.CreateSentence(clipcode)
	if err == nil {
		fmt.Printf("Remember it with: %s\n", mnemonic)
	}
}

func get(clipcode string) {
	data, err := client.Download("http://localhost:8080/", clipcode)
	if err != nil {
		fmt.Printf("Unable to retrieve data: %v\n", err)
		return
	}
	defer data.Close()

	file, err := ioutil.TempFile("", "clippy-*.zip")
	tempFilename := file.Name()
	if err != nil {
		fmt.Printf("Unable to create temporary file: %v\n", err)
		return
	}
	defer os.Remove(tempFilename)

	_, err = io.Copy(file, data)
	file.Close()
	if err != nil {
		fmt.Printf("Unable to save data: %v\n", err)
		return
	}

	err = archiver.Unarchive(tempFilename, ".")
	if err != nil {
		fmt.Printf("Unable to unarchive data: %v\n", err)
		return
	}
}

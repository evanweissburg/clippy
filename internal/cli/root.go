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
		err := put(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case "get":
		clipcode := os.Args[2]
		err := get(clipcode)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	default:
		invalid_usage()
	}
}

func put(filename string) error {
	tempDir, err := ioutil.TempDir("", "clippy-*")
	if err != nil {
		return fmt.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	zipFilename := path.Join(tempDir, "clip.zip")

	err = archiver.Archive([]string{filename}, zipFilename)
	if err != nil {
		return fmt.Errorf("Unable to create zip of %s at %s: %v", filename, zipFilename, err)
	}

	zipFile, err := os.Open(zipFilename)
	if err != nil {
		return fmt.Errorf("Unable to open zip file %s: %v", zipFilename, err)
	}
	defer zipFile.Close()

	clipcode, err := client.Upload("http://localhost:8080/", zipFile)
	if err != nil {
		return fmt.Errorf("Unable to upload to server: %v", err)
	}

	fmt.Printf("Recieved clipcode %s\n", clipcode)

	rand.Seed(time.Now().UTC().UnixNano())
	mnemonic, err := mnemonic.CreateSentence(clipcode)
	if err == nil {
		fmt.Printf("Remember it with: %s\n", mnemonic)
	}

	return nil
}

func get(clipcode string) error {
	data, err := client.Download("http://localhost:8080/", clipcode)
	if err != nil {
		return fmt.Errorf("Unable to retrieve data: %v", err)
	}
	defer data.Close()

	file, err := ioutil.TempFile("", "clippy-*.zip")
	tempFilename := file.Name()
	if err != nil {
		return fmt.Errorf("Unable to create temporary file: %v", err)
	}
	defer os.Remove(tempFilename)

	_, err = io.Copy(file, data)
	file.Close()
	if err != nil {
		return fmt.Errorf("Unable to save data: %v", err)
	}

	err = archiver.Unarchive(tempFilename, ".")
	if err != nil {
		return fmt.Errorf("Unable to unarchive data: %v", err)
	}

	return nil
}

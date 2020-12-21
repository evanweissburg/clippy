package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
	"unicode"
)

var mu sync.Mutex
var db = make(map[string]time.Time)

const (
	serverFileStorageDir = "server_data/"
	clipcodeLength       = 4
)

func Execute() {
	os.RemoveAll(serverFileStorageDir)
	os.Mkdir(serverFileStorageDir, 0700)

	http.HandleFunc("/", handler)
	fmt.Println("Hosting HTTP server on port 8080")
	go makeRefreshTicker(30, 30)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func makeRefreshTicker(refreshSecs time.Duration, keepSecs time.Duration) {
	ticker := time.NewTicker(refreshSecs * time.Second)
	for {
		select {
		case <-ticker.C:
			go enactRefresh(keepSecs)
		}
	}
}

func enactRefresh(keepSecs time.Duration) {
	fmt.Printf("Starting server refresh\n")
	victims := make([]string, 0)
	mu.Lock()
	for clipcode, timestamp := range db {
		if time.Now().Sub(timestamp) > keepSecs*time.Second {
			delete(db, clipcode)
			victims = append(victims, clipcode)
		}
	}
	mu.Unlock()

	for _, clipcode := range victims {
		if err := os.Remove(serverFileStorageDir + clipcode); err != nil {
			fmt.Printf("\tFailed to remove file on server refresh with clipcode %s\n", clipcode)
			continue
		}
	}
	fmt.Printf("\tFinished server refresh, %d objects removed\n", len(victims))
}

func makeClipcode() string {
	for {
		code := make([]rune, clipcodeLength)
		for i := range code {
			code[i] = 'A' + rune(rand.Intn(26))
		}

		if _, ok := db[string(code)]; !ok {
			return string(code)
		}
	}
}

func isClipcode(str string) bool {
	if len(str) != 4 {
		return false
	}
	for _, e := range str {
		if !unicode.IsUpper(e) {
			return false
		}
	}
	return true
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Recieved request %s %s\n", r.Method, r.URL)

	if r.Method == http.MethodPost && r.URL.Path == "/" {
		defer r.Body.Close()
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("\tFailed to decode body\n")
			return
		}

		if len(bytes) == 0 {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Printf("\tEmpty upload attempt\n")
			return
		}

		clipcode := handleUpload(bytes)

		w.Write([]byte(clipcode))
		fmt.Printf("\tUpload successful with clipcode %s\n", clipcode)
	} else if r.Method == http.MethodGet {
		clipcode := path.Base(r.URL.Path)

		if data := handleRequest(clipcode); data != nil {
			w.Write(data)
			fmt.Printf("\tRetrieved valid clipcode %s\n", clipcode)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("\tInvalid clipcode %s\n", clipcode)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Printf("Invalid request\n")
		return
	}
}

func handleUpload(data []byte) string {
	clipcode := makeClipcode()
	err := ioutil.WriteFile(serverFileStorageDir+clipcode, data, 0644)
	if err != nil {
		fmt.Printf("\tFailed to write file with clipcode %s\n", clipcode)
	}

	mu.Lock()
	db[clipcode] = time.Now()
	mu.Unlock()
	return clipcode
}

func handleRequest(clipcode string) []byte {
	if !isClipcode(clipcode) {
		return nil
	}

	data, err := ioutil.ReadFile(serverFileStorageDir + clipcode)
	if err != nil {
		return nil
	}
	return data
}

package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path"
	"sync"
)

var mu sync.Mutex
var db = make(map[string]string)

func Execute() {
	http.HandleFunc("/", handler)
	fmt.Println("Hosting HTTP server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func nextClipcode() string {
	for {
		code := make([]rune, 4)
		for i := range code {
			code[i] = 'a' + rune(rand.Intn(26))
		}

		if _, ok := db[string(code)]; !ok {
			return string(code)
		}
	}
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

		clipcode := handleUpload(string(bytes))

		w.Write([]byte(clipcode))
		fmt.Printf("\tUpload successful with clipcode %s\n", clipcode)
	} else if r.Method == http.MethodGet {
		clipcode := path.Base(r.URL.Path)

		if data, ok := handleRequest(clipcode); ok {
			w.Write([]byte(data))
			fmt.Printf("\tRetrieved valid clipcode %s\n", clipcode)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("\tInvalid clipcode %s\n", clipcode)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Printf("Invalid request\n")
	}
}

func handleUpload(data string) string {
	mu.Lock()
	clipcode := nextClipcode()
	db[clipcode] = data
	mu.Unlock()

	return clipcode
}

func handleRequest(clipcode string) (string, bool) {
	mu.Lock()
	data, ok := db[clipcode]
	mu.Unlock()
	return data, ok
}

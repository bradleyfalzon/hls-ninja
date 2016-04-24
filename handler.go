package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/grafov/m3u8"
)

// Some test identifiers
// shouldn't the tests just register themselves....

//type results struct {
//}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there\n")

	fmt.Fprintf(w, "Go to: /t/test\n") // TODO randomise
}

func results(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fmt.Fprintf(w, "Have player go to: /t/%s.m3u8\n", vars["testID"])

	fmt.Fprintf(w, "==========\n")

	fmt.Fprintf(w, "Test results for: %s\n", vars["testID"])

	// TODO write some auto refresh code

	fmt.Fprintf(w, "")
}

func master(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("Requested: %s", r.URL.String())

	f, err := os.Open(fmt.Sprintf("%s/testsrc.m3u8", VideosDir))
	if err != nil {
		log.Fatal(err) // TODO handle better
	}
	defer f.Close()

	pls := m3u8.NewMasterPlaylist()

	err = pls.DecodeFrom(bufio.NewReader(f), false)
	if err != nil {
		log.Fatal(err) // TODO handle better
	}

	pls.Args = fmt.Sprintf("tid=%s", vars["testID"])

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	fmt.Fprint(w, pls.String())
}

func media(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("Requested: %s", r.URL.String())

	f, err := os.Open(fmt.Sprintf("%s/testsrc/%s", VideosDir, vars["media"]))
	if err != nil {
		log.Fatal(err) // TODO handle better
	}
	defer f.Close()

	pls, err := m3u8.NewMediaPlaylist(1024, 1024)
	if err != nil {
		log.Fatal(err) // TODO handle better
	}

	err = pls.DecodeFrom(bufio.NewReader(f), false)
	if err != nil {
		log.Fatal(err) // TODO handle better
	}

	for _, s := range pls.Segments {
		if s == nil {
			break
		}
		s.URI = s.URI + "?tid=" + vars["testID"]
	}

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	fmt.Fprintf(w, pls.String())
}

func segment(w http.ResponseWriter, r *http.Request) {
	log.Printf("Requested: %s", r.URL.String())

	vars := mux.Vars(r)

	f, err := os.Open(fmt.Sprintf("%s/testsrc/%s/%s", VideosDir, vars["media"], vars["segment"]))
	if err != nil {
		log.Fatal(err) // TODO handle this better
	}

	w.Header().Set("Content-Type", "video/mp2ts") // TODO is this the right header ?

	n, err := io.Copy(w, f)
	if err != nil {
		log.Fatal(err) // TODO handle this better
	}

	log.Printf("Bytes written: %v", n)

}

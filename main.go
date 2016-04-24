package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	VideosDir string = "videos"
)

func main() {
	log.Println("hls-ninja starting")

	r := mux.NewRouter()
	r.HandleFunc("/", homepage)
	r.HandleFunc("/t/{testID}.m3u8", master)
	r.HandleFunc("/t/{testID}/{media}", media)
	r.HandleFunc("/t/{testID}/{media}/{segment}", segment)
	r.HandleFunc("/t/{testID}", results)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":3002", r))
}

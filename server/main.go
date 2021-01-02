package main

import (
	"net/http"
	"time"
)

func serveSlowly(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	time.Sleep(250 * time.Millisecond)
	w.Write([]byte("*yawn*"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/slow", serveSlowly)
	// http.ListenAndServe(":8080", apmhttp.Wrap(mux))
	http.ListenAndServe(":8080", mux)
}

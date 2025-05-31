package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "multiserv: please supply the number of HTTP servers to start")
		os.Exit(1)
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "multiserv: invalid number %v\n", err)
		os.Exit(1)
	}

	if n <= 0 {
		os.Exit(0)
	}

	muxers := newMuxers(n)
	var wg sync.WaitGroup

	for i, mux := range muxers {
		mux.Handle("/", handler(i))
		addr := fmt.Sprintf("localhost:%d", 1024+i)
		log.Printf("starting HTTP server at %s", addr)

		wg.Add(1)
		go func(addr string, mux *http.ServeMux) {
			defer wg.Done()
			if err := http.ListenAndServe(addr, mux); err != nil {
				log.Printf("server at %s exited with error: %v", addr, err)
			}
		}(addr, mux)
	}

	wg.Wait()
	os.Exit(1) // all servers must have returned with error
}

func newMuxers(n int) []*http.ServeMux {
	muxers := make([]*http.ServeMux, n)
	for i := range n {
		muxers[i] = http.NewServeMux()
	}
	return muxers
}

type handler int

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello from handler %d\n", h)
}

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	_ "runtime/pprof"
)

func main() {
	http.HandleFunc("/health", health)
	http.HandleFunc("/", rootHandler)
	addr := ":8080"
	log.Printf("Listening at %v\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(resp http.ResponseWriter, req *http.Request) {
	log.Println("entering root handler")
	user := req.URL.Query().Get("user")
	if user != "" {
		fmt.Fprintf(resp, "hello [%v]\n", user)
	} else {
		fmt.Fprint(resp, "hello [stranger]\n")
	}
	fmt.Fprintf(resp, "========================Deatils of the http request header:=======================================")
	for k, v := range req.Header {
		fmt.Fprintf(resp, "%v = %v\n", k, v)
	}
}

func health(resp http.ResponseWriter, req *http.Request) {
	io.WriteString(resp, "ok")
}

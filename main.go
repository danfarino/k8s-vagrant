package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Printf("Listening.\n")
	log.Panic(http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v", r)
		fmt.Fprintf(w, "hello from kubernetes and vagrant! HOSTNAME=%s\n", os.ExpandEnv("$HOSTNAME"))
	})))
}

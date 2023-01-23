package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func main() {
    path, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }
	fileServer := http.FileServer(http.Dir(path))
    fmt.Printf("Starting test server http://localhost:8080/test/test.html\n")
    http.Handle("/", fileServer)
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
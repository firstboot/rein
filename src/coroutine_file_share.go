package main

import (
	"log"
	"net/http"
)

type coroutineFileShareObj struct {
}

func coroutineFileShare() coroutineFileShareObj {
	return coroutineFileShareObj{}
}

func (obj coroutineFileShareObj) run(port string, sharePath string) {
	http.Handle("/", http.FileServer(http.Dir(sharePath)))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

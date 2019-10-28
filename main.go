package main

import (
	"fmt"
	"git.qietv.work/go-public/health"
	"net/http"
	"time"
)

type helloHandler struct{}

func (h *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func main() {
	http.Handle("/", &helloHandler{})

	//http.ListenAndServe(":8080", nil)

	health.Start(":8080")

	for i := 1; ; {
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}

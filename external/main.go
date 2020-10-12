package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	http.HandleFunc("/", logger(HandleHeavyJob))

	fmt.Println("===== External system is started =====")
	log.Println("listening on :9090")
	http.ListenAndServe(":9090", nil)
}

// HandleHeavyJob sample core proses
func HandleHeavyJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// thread sleep selama 1 detik
	time.Sleep(10 * time.Second)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from external system"))
}

// logger print to console
func logger(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path, r.Method)
		fn(w, r)
	}
}

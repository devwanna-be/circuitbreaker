package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	//hystrix config
	hystrix.ConfigureCommand("command_config", hystrix.CommandConfig{
		Timeout:                1000,
		MaxConcurrentRequests:  300,
		RequestVolumeThreshold: 10,
		SleepWindow:            1000,
		ErrorPercentThreshold:  50,
	})

	http.HandleFunc("/", logger(HandleSubsystem))

	fmt.Println("===== Main system is started =====")
	log.Println("listening on :8080")
	http.ListenAndServe(":8080", nil)
}

// HandleSubsystem send request ke external system
func HandleSubsystem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	resultCh := make(chan []byte)
	errCh := hystrix.Go("command_config", func() error {
		resp, err := http.Get("http://localhost:9090")
		if err != nil {
			return err
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		resultCh <- b
		return nil
	}, nil)

	select {
	case res := <-resultCh:
		log.Println("Request external sistem berhasil:", string(res))
		w.WriteHeader(http.StatusOK)
	case err := <-errCh:
		log.Println("Request external sistem gagal:", err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func logger(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path, r.Method)
		fn(w, r)
	}
}

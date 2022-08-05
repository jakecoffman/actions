package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	go func() {
		time.Sleep(time.Minute)
		log.Fatalln("Times up!")
	}()

	server := &http.Server{
		Addr: "0.0.0.0:8000",
	}
	server.Handler = serve(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Got", string(data))
		w.Write([]byte("world!"))
		go func() {
			time.Sleep(5 * time.Second)
			server.Shutdown(context.Background())
		}()
	})

	log.Println("Serving")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}

type serve func(http.ResponseWriter, *http.Request)

func (s serve) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s(w, r)
}

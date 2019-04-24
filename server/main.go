package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/NBCFB/Octopus"
	"log"
	"net/http"
	"html/template"
	"github.com/iamharvey/HowlerMonkey"
	"time"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t, err := template.ParseFiles("server/index.html")
	if err != nil {
		log.Fatalf("[ERROR] %v", err)

	}

	t.Execute(w, nil)
}

func main() {

	var host string
	var port int

	// Handle command line flag
	flag.StringVar(&host, "h", "localhost", "specify host")
	flag.IntVar(&port, "p", 5678, "specify port")
	flag.Parse()

	// Make up addr with default settings
	addr := fmt.Sprintf("%s:%d", host, port)

	// Create a new broker
	b := HowlerMonkey.NewBroker()

	// Start the broker
	b.Start()

	// Set up router
	r := chi.NewRouter()

	r.Get("/", home)
	r.Get("/events", b.GetEvents)
	r.Get("/send/{event}", b.SendEvent)

	s := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Starts server gracefully
	log.Printf("[INFO] Starting server at http://%s\n", addr)
	Octopus.GracefulServe(s, false)
}
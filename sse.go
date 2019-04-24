package HowlerMonkey

import (
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"time"
)

/*
Broker's behaviours can be defined as follows:
	- add a new connection as a client;
	- remove a connected client;
	- broadcast event to currently connected clients;
	- sync.WaitGroup ensures the background thread is finished before shutdown is completed.
 */
type Broker struct {
	// Current connected clients
	clients			map[chan []byte]int

	// Add new clients
	newClients 		chan chan []byte

	// Remove clients
	removeClients 	chan chan []byte

	// Message to be broadcast to connected clients
	Events 			chan []byte
}

/*
NewBroker creates a broker instance.
 */
func NewBroker() (b *Broker) {
	b = &Broker{
		clients:		make(map[chan []byte]int),
		newClients:		make(chan (chan []byte)),
		removeClients:	make(chan (chan []byte)),
		Events:			make(chan []byte),
	}
	return
}

//Start starts the server
func (b *Broker) Start() {
	go b.listen()
}

/*
Listen listens to the signals from different channels and act accordingly:
- If new signal is detected from newClients, add a new client;
- If new signal is detected from removeClients, remove a client;
- If new signal is detected from Events, broadcast received events to all the connected clients.
 */
func (b *Broker) listen() {
	for {
		numClients := len(b.clients)
		select {
		case c := <-b.newClients:
			b.clients[c] = time.Now().Nanosecond()
			log.Printf("[INFO] Add a new client. Current num-of-client: %d.\n", numClients)

		case c := <-b.removeClients:
			delete(b.clients, c)
			close(c)
			log.Printf("[INFO] Remove a client. Current num-of-client: %d.\n", numClients)

		case e := <-b.Events:
			for c := range b.clients {
				c <- e
			}
			log.Printf("[INFO] Broadcasted %q to %d clients.\n", e, numClients)
		}
	}
}

/*
GetEvents is a http handler func that expects and receives events from the server.
 */
func (b *Broker) GetEvents(w http.ResponseWriter, r *http.Request) {

	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Create a new channel for sending event to the client.
	eChan := make(chan []byte)

	b.newClients <- eChan

	// Listen to the closing of the http connection via the CloseNotifier
	notify := r.Context().Done()
	go func() {
		<-notify
		// Remove this client from the map of attached clients
		// when `EventHandler` exits.
		b.removeClients <- eChan
		log.Println("HTTP connection just closed.")
	}()

	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	for {

		e, open := <-eChan

		if !open {
			break
		}

		fmt.Fprintf(w, "data: Event: %s\n\n", e)

		f.Flush()
	}

}

/*
SendEvent allows a client to send an event for broadcasting (via GET method)
*/
func (b *Broker) SendEvent(w http.ResponseWriter, r *http.Request) {
	e := chi.URLParam(r, "event")
	b.Events <- []byte(e)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Event sent: %s\n\n", e)
}




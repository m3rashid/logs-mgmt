package utils

import (
	"fmt"
	"log"
	"net/http"
)

type StreamBroker struct {
	Notifier       chan []byte
	newClients     chan chan []byte
	closingClients chan chan []byte
	clients        map[chan []byte]bool
}

func (broker *StreamBroker) GetClientCount() int {
	return len(broker.clients)
}

func (broker *StreamBroker) listen() {
	for {
		select {
		case s := <-broker.newClients: // add the new connected client
			broker.clients[s] = true
			log.Printf("Client added. %d registered clients", len(broker.clients))

		case s := <-broker.closingClients: // remove the disconnected client
			delete(broker.clients, s)
			log.Printf("Removed client. %d registered clients", len(broker.clients))

		case event := <-broker.Notifier: // got a new event, propagate to all
			for clientMessageChan := range broker.clients {
				clientMessageChan <- event
			}
		}
	}
}

func (broker *StreamBroker) ServeRequestHttp(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	messageChan := make(chan []byte)

	broker.newClients <- messageChan // signal the broker about the new client

	defer func() {
		broker.closingClients <- messageChan // remove the client when it exits
	}()

	notify := rw.(http.CloseNotifier).CloseNotify() // listen to connection close
	go func() {
		<-notify
		broker.closingClients <- messageChan
	}()

	// block till we have a message
	for {
		fmt.Fprintf(rw, "data: %s\n\n", <-messageChan)
		flusher.Flush()
	}
}

func NewStreamBroker() *StreamBroker {
	broker := &StreamBroker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	go broker.listen()

	return broker
}

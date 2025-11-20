package stream

import "log"

type Broker struct {
	Notifier       chan []byte
	newClients     chan chan []byte
	closingClients chan chan []byte
	clients        map[chan []byte]bool
}

func (broker *Broker) listen() {
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

func NewBroker() (broker *Broker) {
	broker = &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	go broker.listen()

	return
}

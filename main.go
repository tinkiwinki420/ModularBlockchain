package main

import (
	"modularBlockchain/network" // Importing your custom network package.
	"time"                      // Importing time for adding delays.
)

func main() {
	// Create two LocalTransport instances: one representing "LOCAL" and the other "REMOTE".
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	// Establish a bi-directional connection between the two transports.
	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	// Start a goroutine that will continuously send messages from "REMOTE" to "LOCAL".
	go func() {
		for {
			// Send a message from "REMOTE" to "LOCAL" every second.
			trRemote.SendMessage(trLocal.Addr(), []byte("hello mfers"))
			time.Sleep(1 * time.Second) // Pause for 1 second before sending the next message.
		}
	}()

	// Define server options, setting "LOCAL" transport to be managed by the server.
	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal},
	}

	// Create a new server with the specified options.
	s := network.NewServer(opts)

	// Start the server, which will begin consuming RPCs and handling messages.
	s.Start()
}

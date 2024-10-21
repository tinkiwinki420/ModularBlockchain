package network

import (
	"fmt"
	"time"
)

// ServerOpts holds configuration options for creating a Server.
// It includes a list of Transports that the server will manage.
type ServerOpts struct {
	Transports []Transport // List of transports that this server will use for communication.
}

// Server represents a network server that manages incoming RPCs (Remote Procedure Calls)
// and handles them accordingly.
type Server struct {
	ServerOpts               // Embeds ServerOpts to access configuration directly.
	rpcCh      chan RPC      // Channel to receive RPCs from different transports.
	quitCh     chan struct{} // Channel to signal when to shut down the server.
}

// NewServer creates a new Server instance with the provided options.
// It initializes channels for handling RPCs and server shutdown.
func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		rpcCh:      make(chan RPC),         // Creates an unbuffered channel for receiving RPCs.
		quitCh:     make(chan struct{}, 1), // Creates a buffered channel for shutdown signaling.
	}
}

// Start begins the server's main loop, processing RPCs and handling periodic tasks.
// It reads from the rpcCh for incoming messages and the quitCh for shutdown signals.
func (s *Server) Start() {
	// Initialize all transports so they can start consuming messages.
	s.initTransports()

	// Create a ticker that triggers every 5 seconds.
	ticker := time.NewTicker(5 * time.Second)

free:
	for {
		select {
		// If an RPC is received from the rpcCh, print its details.
		case rpc := <-s.rpcCh:
			fmt.Printf("%+v\n", rpc)

		// If a message is received on quitCh, break out of the loop to stop the server.
		case <-s.quitCh:
			break free

		// If the ticker triggers (every 5 seconds), print "meow".
		case <-ticker.C:
			fmt.Println("meow")
		}
	}

	// When the loop ends (on quit signal), print a shutdown message.
	fmt.Println("server shutdown")
}

// initTransports starts a goroutine for each transport to consume RPCs.
// It forwards these RPCs to the server's rpcCh for centralized processing.
func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			// Continuously consume messages from the transport.
			for rpc := range tr.Consume() {
				// Forward each received RPC to the server's main rpcCh.
				s.rpcCh <- rpc
			}
		}(tr)
	}
}

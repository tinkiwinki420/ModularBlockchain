package network

import (
	"fmt"
	"sync"
)

// LocalTransport simulates a transport mechanism in a network, allowing nodes to communicate.
// It holds information about its address, communication channel, connected peers, and a lock for concurrency.
type LocalTransport struct {
	addr      NetAddr                     // Address (identifier) of this transport.
	consumeCh chan RPC                    // Channel to receive incoming RPC messages.
	lock      sync.RWMutex                // Read-write lock to ensure thread-safe access to peers.
	peers     map[NetAddr]*LocalTransport // Map of connected peers, using their addresses as keys.
}

// NewLocalTransport creates a new LocalTransport instance with the specified address.
// Returns it as a Transport interface, allowing for potential flexibility in using different implementations.
func NewLocalTransport(addr NetAddr) Transport {
	return &LocalTransport{
		addr:      addr,
		consumeCh: make(chan RPC, 1024),              // Buffered channel to hold incoming messages.
		peers:     make(map[NetAddr]*LocalTransport), // Initialize an empty map for peers.
	}
}

// Consume returns the consumeCh channel, allowing other components to receive messages from this transport.
func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeCh
}

// Connect establishes a connection with another Transport by adding it to the peers map.
// This method only allows connections to other LocalTransports, hence the type assertion.
func (t *LocalTransport) Connect(tr Transport) error {
	t.lock.Lock()         // Acquire write lock to prevent concurrent modifications of peers.
	defer t.lock.Unlock() // Ensure the lock is released after the operation.

	// Add the transport to the peers map, using its address as the key.
	t.peers[tr.Addr()] = tr.(*LocalTransport)

	return nil // Return nil indicating the connection was successful.
}

// SendMessage sends a message (payload) to a specified peer address.
// The message is wrapped in an RPC struct and sent to the peer's consume channel.
func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.lock.RLock()         // Acquire read lock to ensure the peers map is accessed safely.
	defer t.lock.RUnlock() // Release the lock after reading.

	// Look up the peer in the peers map using the provided address.
	peer, ok := t.peers[to]
	if !ok {
		// If the peer is not found, return an error indicating the failure.
		return fmt.Errorf("%s: couldn't send message to %s", t.addr, to)
	}

	// Send the message as an RPC to the peer's consume channel.
	peer.consumeCh <- RPC{
		From:    t.addr,  // Sender's address.
		Payload: payload, // The message content.
	}
	return nil // Return nil indicating the message was successfully sent.
}

// Addr returns the address of this LocalTransport.
func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}

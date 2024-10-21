// local_transport_test.go
package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	err1 := tra.Connect(trb)
	err2 := trb.Connect(tra)

	assert.NoError(t, err1)
	assert.NoError(t, err2)

	assert.Equal(t, tra.peers[trb.Addr()], trb)
	assert.Equal(t, trb.peers[tra.Addr()], tra)
}
func TestSendMsg(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	msg := []byte("hello nigger")
	assert.Nil(t, tra.SendMessage(trb.addr, msg))

	rpc := <-trb.Consume()
	assert.Equal(t, rpc.Payload, msg)
	assert.Equal(t, rpc.From, tra.addr)

}

package pkg

import (
	"testing"
	"time"
)

func TestListenBroadcast(t *testing.T) {
	TestListener := make(chan MessageResponse)
	for len(TestListener) > 0 {
		<-TestListener
	}
	go Listen(TestListener)

	time := time.Now()
	message := "test"
	messageType := "test"

	m := Message{Type: messageType, Message: message, Timestamp: time}
	Broadcast(m)

	response := <-TestListener

	if response.Message.Type != messageType {
		t.Errorf("MessageType was incorrect, got: %s, want: %s.", response.Message.Type, messageType)
	}

	if response.Message.Message != message {
		t.Errorf("Message was incorrect, got: %s, want: %s.", response.Message.Message, message)
	}

	if !response.Message.Timestamp.Equal(time) {
		t.Errorf("Timestamp was incorrect, got: %s, want: %s.", response.Message.Timestamp, time)
	}

}

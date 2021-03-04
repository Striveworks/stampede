package pkg

import (
	"net"
	"time"
)

// Message represents a message from one node to another
type Message struct {
	Type      string    `json:"type"`      // Type of message
	Timestamp time.Time `json:"time"`      // Timestamp of message creation
	Recipient string    `json:"recipient"` // Desired recipient of message
	Message   string    `json:"message"`   // Message body
	Node      Node      `json:"node"`      // Originator node
}

// MessageResponse represents a potential response message
type MessageResponse struct {
	Message    Message        `json:"message"`    // Message struct
	Address    net.Addr       `json:"address"`    // Network address
	Connection net.PacketConn `json:"connection"` // Network connection
}

// LeaderAsk sends election vote over multicast
func LeaderAsk(node Node) {
	m := Message{Type: "Election", Message: "Vote", Timestamp: time.Now(), Node: node}
	Broadcast(m)
}

// HeartBeat sends heartbeat over multicast
func HeartBeat(node Node) {
	m := Message{Type: "Heartbeat", Message: "Alive", Timestamp: time.Now(), Node: node}
	Broadcast(m)
}

// JoinRequest sends requests to join the cluster over multicast
func JoinRequest(node Node) {
	m := Message{Type: "JoinRequest", Message: "", Timestamp: time.Now(), Node: node}
	Broadcast(m)
}

// JoinResponse sends join keys to a specified node over multicast
func JoinResponse(uuid string, creds string) {
	m := Message{Type: "JoinResponse", Recipient: uuid, Message: creds, Timestamp: time.Now()}
	Broadcast(m)
}

package pkg

import (
	"net"
	"strings"
	"time"
)

type Message struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"time"`
	Recipient string    `json:"recipient"`
	Message   string    `json:"message"`
	Node      Node      `json:"node"`
}

type MessageResponse struct {
	Message    Message        `json:"message"`
	Address    net.Addr       `json:"address"`
	Connection net.PacketConn `json:"connection"`
}

func LeaderAsk(node Node) {
	m := Message{Type: "Election", Message: "Vote", Timestamp: time.Now(), Node: node}
	Broadcast(m)
}

func HeartBeat(node Node) {
	m := Message{Type: "Heartbeat", Message: "Alive", Timestamp: time.Now(), Node: node}
	Broadcast(m)
}

func JoinRequest(node Node) {
	m := Message{Type: "JoinRequest", Message: "", Timestamp: time.Now(), Node: node}
	Broadcast(m)
}

func JoinResponse(uuid string, keys []string) {
	m := Message{Type: "JoinResponse", Recipient: uuid, Message: strings.Join(keys, " "), Timestamp: time.Now()}
	Broadcast(m)
}

package pkg

import (
	"net"
	"time"
)

type LeaderMessage struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"time"`
	Message   string    `json:"message"`
	Node      Node      `json:"node"`
}

type LeaderResponse struct {
	LeaderMessage LeaderMessage  `json:"leader"`
	Address       net.Addr       `json:"address"`
	Connection    net.PacketConn `json:"connection"`
}

func BecomeLeader(node Node) {
	m := LeaderMessage{Type: "Election", Message: "Win", Timestamp: time.Now(), Node: node}
	Broadcast(m)
}

func LeaderAsk(node Node) {
	m := LeaderMessage{Type: "Election", Message: "Vote", Timestamp: time.Now(), Node: node}
	Broadcast(m)
}

func DenyElection(node Node) {
	m := LeaderMessage{Type: "Election", Message: "Denied", Timestamp: time.Now(), Node: node}
	Broadcast(m)

}

func HeartBeat(node Node) {
	m := LeaderMessage{Type: "Heartbeat", Message: "Alive", Timestamp: time.Now(), Node: node}
	Broadcast(m)
}

func LeaderEnforce(node Node) {
	m := LeaderMessage{Type: "Heartbeat", Message: "Leader", Timestamp: time.Now(), Node: node}
	Broadcast(m)
}

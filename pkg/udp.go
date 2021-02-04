package pkg

import (
	"encoding/json"
	"log"
	"net"
)

const (
	address = "224.0.0.1:9999"
)

//Broadcast ..
func Broadcast(leaderMessage LeaderMessage) {
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp4", addr)

	var jsonData []byte
	jsonData, err = json.Marshal(leaderMessage)
	if err != nil {
		log.Println(err)
	}
	_, err = conn.WriteToUDP(jsonData, addr)
	if err != nil {
		panic(err)
	}

}

//Listen ..
func Listen(c chan LeaderResponse) {
	defer close(c)
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		panic(err)
	}
	// Loop forever reading from the socket
	for {
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}
		data := buf[:n]
		var response LeaderMessage
		err = json.Unmarshal(data, &response)

		result := LeaderResponse{LeaderMessage: response, Address: addr, Connection: conn}

		c <- result

	}
}

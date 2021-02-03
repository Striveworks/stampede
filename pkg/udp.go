package pkg

import (
	"encoding/json"
	"log"
	"net"
)

const (
	sendAddress   = "127.0.0.1:9999"
	listenAddress = ":9999"
)

//Broadcast ..
func Broadcast(leaderMessage LeaderMessage) {
	addr, err := net.ResolveUDPAddr("udp4", sendAddress)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp4", nil, addr)

	var jsonData []byte
	jsonData, err = json.Marshal(leaderMessage)
	if err != nil {
		log.Println(err)
	}
	_, err = conn.Write(jsonData)
	if err != nil {
		panic(err)
	}

}

//Listen ..
func Listen(c chan LeaderResponse) {
	defer close(c)
	conn, err := net.ListenPacket("udp4", listenAddress)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

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

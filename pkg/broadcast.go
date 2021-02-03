package pkg

import (
  "net"
  "fmt"
)

func Broadcast() {
  pc, err := net.ListenPacket("udp4", ":8829")
  if err != nil {
    panic(err)
  }
  defer pc.Close()

  addr,err := net.ResolveUDPAddr("udp4", "192.168.7.255:8829")
  if err != nil {
    panic(err)
  }

  _,err := pc.WriteTo([]byte("data to transmit"), addr)
  if err != nil {
    panic(err)
  }
}

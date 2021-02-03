package pkg

import (
  "net"
  "fmt"
)


func Listen() {
  pc,err := net.ListenPacket("udp4", ":8829")
  if err != nil {
    panic(err)
  }
  defer pc.Close()

  buf := make([]byte, 1024)
  n,addr,err := pc.ReadFrom(buf)
  if err != nil {
    panic(err)
  }

  fmt.Printf("%s sent this: %s\n", addr, buf[:n])
}

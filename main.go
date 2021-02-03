package main

import (
	"striveworks.us/stampede/pkg"
)

const (
	sendAddress   = "127.0.0.1:9999"
	listenAddress = ":9999"
)

func main() {
	node := pkg.New(&pkg.NodeConfig{IsLeader: false})
	node.Start()

}

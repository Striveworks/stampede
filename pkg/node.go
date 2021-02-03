package pkg

import (
	"time"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	waitInterval     = 3
	electionInterval = 3
)

type Node struct {
	UUID         string    `json:"uuid"`
	IsLeader     bool      `json:"isleader"`
	Voting       bool      `json:"voting"`
	ElectionTime time.Time `json:"election"`
}

//NodeConfig ...
type NodeConfig struct {
	IsLeader bool
}

//New ...
func New(cfg *NodeConfig) Node {

	return Node{
		UUID:     generateUUID().String(),
		IsLeader: cfg.IsLeader,
	}
}

func (node Node) Start() error {
	listener := make(chan LeaderResponse)
	//Use a nodePool map so lookup times are O(1)
	nodePool := make(map[string]Node)
	electionDenials := 0
	startTime := time.Now()

	go Listen(listener)

	for {
		// Create non-blocking channel to listen for UDP messages
		select {
		case response := <-listener:
			log.Info(response)
			//Only care about messages that aren't from myself
			if response.LeaderMessage.Node.UUID != node.UUID {

				//Add other nodes to NodePool if they do not exist
				if _, ok := nodePool[node.UUID]; ok {
				} else {
					nodePool[node.UUID] = node
				}

				//Increment Election denial
				if response.LeaderMessage.Type == "Election" && response.LeaderMessage.Message == "Denied" {
					electionDenials++
				}

				//Reject other leader elections if I am the leader
				if node.IsLeader && response.LeaderMessage.Type == "Election" && response.LeaderMessage.Message == "Vote" {
					DenyElection(node)
				}

				//TODO
				//If I get heartbeats from other nodes and I am the leader, send back shared key

			}
		default:
		}

		waitElapsed := time.Since(startTime).Seconds()
		electionElapsed := time.Since(node.ElectionTime).Seconds()

		//If I am voting and have not heard any denies in election interval, become Leader
		if !node.IsLeader && node.Voting && electionElapsed > electionInterval && electionDenials == 0 {
			BecomeLeader(node)
			node.IsLeader = true
			log.Info("Have heard ", electionDenials, " denials")
			log.Info("Assuming leader role")
		}

		//If I haven't gotten any heartbeats from other nodes in specified interval and I haven't already voted, make leader election
		if !node.IsLeader && waitElapsed > waitInterval && !node.Voting && len(nodePool) == 0 {
			node.Voting = true
			node.ElectionTime = time.Now()
			LeaderAsk(node)
			log.Info("Sending election request")
		}

		//TODO
		//Case when multiple nodes are competing for leader role
		HeartBeat(node)
		time.Sleep(1 * time.Second)
	}
}

var nsUUID = uuid.Must(uuid.FromString("34b13033-50e7-4083-97f5-d389cf3a1c0e"))

func generateUUID() uuid.UUID {
	id, err := uuid.NewV1()
	if err != nil {
		id, err = uuid.NewV4()
		if err != nil {
			return uuid.NewV5(nsUUID, time.Now().String())
		}
	}

	return id
}

package pkg

import (
	"time"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	voteCount = 10
)

type Node struct {
	UUID          string    `json:"uuid"`
	IsLeader      bool      `json:"isleader"`
	Voting        bool      `json:"voting"`
	ElectionTime  time.Time `json:"election"`
	LastHeartBeat time.Time `json:"hearbeat"`
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
	// electionDenials := 0
	votes := 0

	go Listen(listener)

	for {
		// Create non-blocking channel to listen for UDP messages
		select {
		case response := <-listener:
			//Only care about messages that aren't from myself
			if response.LeaderMessage.Node.UUID != node.UUID {
				log.Info(response)
				//Reset Node Heartbeat time
				response.LeaderMessage.Node.LastHeartBeat = time.Now()

				//Add other nodes to NodePool
				nodePool[response.LeaderMessage.Node.UUID] = response.LeaderMessage.Node

				// //Increment Election denial
				// if response.LeaderMessage.Type == "Election" && response.LeaderMessage.Message == "Denied" {
				// 	electionDenials++
				// }
				//
				// //Reject other leader elections if I am the leader
				// if node.IsLeader && response.LeaderMessage.Type == "Election" && response.LeaderMessage.Message == "Vote" {
				// 	log.Info("Blocking the election")
				// 	DenyElection(node)
				// }

				//TODO
				//If I get heartbeats from other nodes and I am the leader, send back shared key

			}
		default:
		}

		//If I am voting and have not heard any denies in election period, become Leader
		if !node.IsLeader && node.Voting && votes >= voteCount {
			BecomeLeader(node)
			node.IsLeader = true
			// log.Info("Have heard ", electionDenials, " denials")
			log.Info("Assuming leader role")
		}

		//If I haven't already voted and it has been longer than the waitInterval, make leader election
		if !node.IsLeader {

			//Set intial vote time
			if !node.Voting {
				node.ElectionTime = time.Now()
			}

			//If I haven't gotten any heartbeats from other nodes
			if len(nodePool) == 0 {
				node.Voting = true
				LeaderAsk(node)
				votes++
				log.Info(votes, "/", voteCount, " votes")
			} else {
				//If I have gotten heartbeats, but the ALL other node's election
				// time came after mine still send a vote
				latest, err := latestElection(nodePool)
				if err != nil {
					log.Error(err)
				}
				log.Info(latest.UUID)
				if latest.ElectionTime.After(node.ElectionTime) {
					node.Voting = true
					LeaderAsk(node)
					votes++
					log.Info(votes, "/", voteCount, " votes")
				}

			}
		}

		HeartBeat(node)

		log.Info(len(nodePool))
		for _, v := range nodePool {
			if time.Since(v.LastHeartBeat).Seconds() > 10 {
				delete(nodePool, v.UUID)
				log.Info("Deleted ", v.UUID, " from nodes")
			}
		}

		if node.IsLeader {
			log.Info("I am the captain now!")
		} else {
			log.Info("Following")
		}
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

func latestElection(nodePool map[string]Node) (Node, error) {
	latest := Node{ElectionTime: time.Now()}
	for _, v := range nodePool {
		if !v.ElectionTime.IsZero() && v.ElectionTime.Before(latest.ElectionTime) {
			latest = v
		}
	}
	return latest, nil
}

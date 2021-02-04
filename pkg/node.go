package pkg

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	voteCount = 10
	stateFile = "/opt/stampede/is-joined"
)

var (
	nodePool    map[string]Node
	currentNode Node
	votes       int
)

type Node struct {
	UUID            string    `json:"uuid"`
	IsLeader        bool      `json:"isleader"`
	IsJoined        bool      `json:"isjoined"`
	LastJoinRequest time.Time `json:"lastjoinrequest"`
	Voting          bool      `json:"voting"`
	ElectionTime    time.Time `json:"election"`
	LastHeartBeat   time.Time `json:"hearbeat"`
}

//CreateNode ...
func CreateNode() Node {
	return Node{
		UUID:            generateUUID().String(),
		IsLeader:        false,
		LastJoinRequest: time.Now(),
	}
}

func (node Node) Start() {
	if _, err := os.Stat(stateFile); err == nil {
		log.Info("Already a cluster member")
		os.Exit(0)
	}

	currentNode = node
	nodePool = make(map[string]Node)
	votes = 0

	go recieve()

	for {
		go cleanNodePool()

		if !currentNode.IsLeader && currentNode.Voting && votes >= voteCount {
			currentNode.IsLeader = true
			log.Info("I am the captain now!")

		}

		if !currentNode.IsLeader {
			log.Info("Following")

			if !currentNode.Voting {
				currentNode.ElectionTime = time.Now()
			}
			electNode()

			if time.Since(currentNode.LastJoinRequest).Seconds() > 20 && !currentNode.IsJoined {
				currentNode.LastJoinRequest = time.Now()
				JoinRequest(currentNode)
			}
		}
		HeartBeat(currentNode)

		time.Sleep(1 * time.Second)
	}
}

func electNode() {
	latest, err := latestElection(nodePool)
	if err != nil {
		log.Error(err)
	}
	if len(nodePool) == 0 || latest.ElectionTime.After(currentNode.ElectionTime) {
		currentNode.Voting = true
		LeaderAsk(currentNode)
		votes++
		log.Info(votes, "/", voteCount, " votes")
	}
}

func cleanNodePool() {
	for _, v := range nodePool {
		if time.Since(v.LastHeartBeat).Seconds() > 30 {
			delete(nodePool, v.UUID)
			log.Info("Deleted ", v.UUID, " from nodes")
		}
	}
}

func recieve() {
	listener := make(chan MessageResponse)
	go Listen(listener)

	// Create non-blocking channel to listen for UDP messages
	for {
		select {
		case response := <-listener:
			if response.Message.Node.UUID != currentNode.UUID {
				response.Message.Node.LastHeartBeat = time.Now()

				nodePool[response.Message.Node.UUID] = response.Message.Node

				if response.Message.Type == "JoinRequest" && currentNode.IsLeader {
					addNode(response)
				}

				if response.Message.Type == "JoinResponse" && response.Message.Recipient == currentNode.UUID {
					joinCluster(response)
				}
			}
		default:
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func addNode(response MessageResponse) {
	app := "microk8s"
	arg := "add-node"

	cmd := exec.Command(app, arg)
	stdout, err := cmd.Output()
	if err != nil {
		log.Error(err)
	}
	msg := strings.Split(string(stdout), "\n")
	log.Info(msg)
	log.Info("Allowing ", response.Address, ": ", response.Message.Node.UUID, " to join. Sending keys")
	JoinResponse(response.Message.Node.UUID, msg[len(msg)-5:])
}

func joinCluster(response MessageResponse) {
	for _, key := range strings.Split(response.Message.Message, " microk8s join ") {
		app := "microk8s"
		arg := "join"

		cmd := exec.Command(app, arg, key)
		_, err := cmd.Output()
		if err == nil {
			currentNode.IsJoined = true
			_, err := os.Create(stateFile)
			if err != nil {
				log.Fatal(err)
			}
			log.Info("Joined cluster, shutting down...")
			os.Exit(0)
		}
	}
	if !currentNode.IsJoined {
		log.Error("Failed to join cluster")
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

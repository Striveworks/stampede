package pkg

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	voteCount = 10
	stateFile = "/opt/stampede-is-joined"
	uuidNs    = "34b13033-50e7-4083-97f5-d389cf3a1c0e"
)

var (
	nodePool    map[string]Node
	currentNode Node
	votes       int
)

// Node represents each individual node instance
type Node struct {
	UUID            string    `json:"uuid"`            // UUID to uniquely identify each node
	IsLeader        bool      `json:"isleader"`        // Identifies if node is leader
	IsJoined        bool      `json:"isjoined"`        // Identifies if node has joined the cluster
	LastJoinRequest time.Time `json:"lastjoinrequest"` // Timestamp of last join request
	Voting          bool      `json:"voting"`          // Identifies if node is a voting member
	ElectionTime    time.Time `json:"election"`        // Timestamp of initial election action
	LastHeartBeat   time.Time `json:"hearbeat"`        // Timestamp of last hearbeat
	ClusterCreds    string    `json:"clustercreds"`
}

// CreateNode returns a node instance with defaults set
func CreateNode() Node {
	return Node{
		UUID:            generateUUID().String(),
		IsLeader:        false,
		LastJoinRequest: time.Now(),
	}
}

// Start runs the main node loop. It will run until the node has become the
// leader or the node has joined the cluster
func (node Node) Start() {
	if _, err := os.Stat(stateFile); err == nil {
		log.Info("Already a cluster member")
		os.Exit(0)
	}

	currentNode = node
	nodePool = make(map[string]Node)
	votes = 0

	go receive()

	for {
		go cleanNodePool(nodePool)

		if !currentNode.IsLeader && currentNode.Voting && votes >= voteCount {
			currentNode.IsLeader = true
			log.Info("I am the captain now!")
			go initCluster()

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

// Begins leader election
func electNode() {
	earliest, err := earliestElection(nodePool)
	if err != nil {
		log.Error(err)
	}
	if len(nodePool) == 0 || earliest.ElectionTime.After(currentNode.ElectionTime) {
		currentNode.Voting = true
		LeaderAsk(currentNode)
		votes++
		log.Info(votes, "/", voteCount, " votes")
	}
}

// Removes "dead" nodes
func cleanNodePool(nodePool map[string]Node) {
	for _, v := range nodePool {
		if time.Since(v.LastHeartBeat).Seconds() > 30 {
			delete(nodePool, v.UUID)
			log.Info("Deleted ", v.UUID, " from nodes")
		}
	}
}

// Handles message processing
func receive() {
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

func initCluster() {
	if viper.GetString("cluster-type") == "kubeadm" {

		var cmd *exec.Cmd
		if viper.GetString("advertise-address") != "" {
			cmd = exec.Command("kubeadm", "init", "--apiserver-advertise-address", viper.GetString("advertise-address"))
		} else {
			cmd = exec.Command("kubeadm", "init")
		}

		stdout, err := cmd.Output()
		if err != nil {
			log.Error(err)
		}
		msg := strings.Split(string(stdout), "\n")
		creds := strings.ReplaceAll(strings.Join(msg[len(msg)-3:], " "), "\\", "")
		log.Info(creds)
		currentNode.ClusterCreds = creds
	}
}

// Adds node to microk8s cluster
func addNode(response MessageResponse) {
	log.Info("Allowing ", response.Address, ": ", response.Message.Node.UUID, " to join. Sending creds")
	if viper.GetString("cluster-type") == "microk8s" {
		addNodeMicroK8s(response)
	}
	if viper.GetString("cluster-type") == "kubeadm" {
		JoinResponse(response.Message.Node.UUID, currentNode.ClusterCreds)
	}
}

// Adds node to microk8s cluster
func addNodeMicroK8s(response MessageResponse) {
	app := "microk8s"
	arg := "add-node"

	cmd := exec.Command(app, arg)
	stdout, err := cmd.Output()
	if err != nil {
		log.Error(err)
	}
	msg := strings.Split(string(stdout), "\n")
	JoinResponse(response.Message.Node.UUID, strings.Join(msg[len(msg)-5:], " "))
}

// Joins node to microk8s cluster
func joinCluster(response MessageResponse) {
	if viper.GetString("cluster-type") == "microk8s" {
		joinClusterMicroK8s(response)
	}
	if viper.GetString("cluster-type") == "kubeadm" {
		joinClusterKubeadm(response)
	}
}

// Joins node to microk8s cluster
func joinClusterMicroK8s(response MessageResponse) {
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

// Joins node to kubeadm cluster
func joinClusterKubeadm(response MessageResponse) {
	ip := strings.Join(strings.Fields(response.Message.Message)[2:3], "")
	token := strings.Join(strings.Fields(response.Message.Message)[4:5], "")
	hash := strings.Join(strings.Fields(response.Message.Message)[6:], "")
	cmd := exec.Command("kubeadm", "join", ip, "--token", token, "--discovery-token-ca-cert-hash", hash)
	_, err := cmd.Output()
	if err == nil {
		currentNode.IsJoined = true
		_, err := os.Create(stateFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("Joined cluster, shutting down...")
		os.Exit(0)
	} else {
		log.Error(err)
	}

	if !currentNode.IsJoined {
		log.Error("Failed to join cluster")
	}
}

// Generates unique identifier
func generateUUID() uuid.UUID {
	nsUUID := uuid.Must(uuid.FromString(uuidNs))
	id, err := uuid.NewV1()
	if err != nil {
		id, err = uuid.NewV4()
		if err != nil {
			return uuid.NewV5(nsUUID, time.Now().String())
		}
	}
	return id
}

// Finds the node with the earliest election time
func earliestElection(nodePool map[string]Node) (Node, error) {
	earliest := Node{ElectionTime: time.Now()}
	for _, v := range nodePool {
		if !v.ElectionTime.IsZero() && v.ElectionTime.Before(earliest.ElectionTime) {
			earliest = v
		}
	}
	return earliest, nil
}

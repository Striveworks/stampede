package pkg

import (
	"testing"
	"time"
)

func TestNode(t *testing.T) {
	node := CreateNode()
	if node.IsLeader != false {
		t.Errorf("Bad node intialization, got: %v, want: %v.", node.IsLeader, false)
	}

}

func TestEarliestElection(t *testing.T) {
	nodePool := make(map[string]Node)
	var firstNode Node

	for i := 0; i < 10; i++ {
		node := CreateNode()
		node.ElectionTime = time.Now()
		if i == 0 {
			firstNode = node
		}
		nodePool[node.UUID] = node
	}

	earliest, err := earliestElection(nodePool)
	if err != nil {
		t.Error(err)
	}
	if earliest != firstNode {
		t.Errorf("Not the earliest node got: %v, want: %v.", earliest, firstNode)
	}

}

func TestCleanNodePool(t *testing.T) {
	nodePool := make(map[string]Node)
	node := CreateNode()
	nodePool[node.UUID] = node
	cleanNodePool(nodePool)
	if len(nodePool) > 0 {
		t.Error("Failed to clean 'dead' nodes")
	}

}

func TestGenerateUUID(t *testing.T) {
	uuid := generateUUID()
	if len(uuid.String()) < 1 {
		t.Error("Failed to generate UUID")
	}
}

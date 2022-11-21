package oddslib

import (
	"math/rand"
	"testing"

	"github.com/MaximeWeyl/fake"
	"github.com/stretchr/testify/assert"
)

// TestInsertQueue : testing that insertions into the priority queue
// keeps the queue ordered, and that removals works properly
func TestInsertRemoveQueue(t *testing.T) {
	queue := DijkstraPriorityQueue{}

	earth := DijkstraNode{
		index: PlanetIndex{
			PlanetName: "Earth",
			FuelLeft:   10,
			Time:       0,
		},
		Cost:             10,
		PreviousPathNode: nil,
		NextQueueNode:    nil,
	}
	err := queue.Insert(&earth) // Insert at first position
	assert.NoError(t, err)
	assert.True(t, &earth == queue.firstItem, "Should have same addresses")
	assert.Nil(t, earth.NextQueueNode)

	mars := DijkstraNode{
		index: PlanetIndex{
			PlanetName: "Mars",
			FuelLeft:   20,
			Time:       1,
		},
		Cost:             5,
		PreviousPathNode: nil,
		NextQueueNode:    nil,
	}
	err = queue.Insert(&mars) // Insert at first position
	assert.NoError(t, err)
	assert.Nil(t, earth.NextQueueNode)
	assert.NotNil(t, mars.NextQueueNode)
	assert.True(t, &mars == queue.firstItem, "Should have same addresses")
	assert.True(t, &earth == queue.firstItem.NextQueueNode, "Should have same addresses")

	saturn := DijkstraNode{
		index: PlanetIndex{
			PlanetName: "Saturn",
			FuelLeft:   50,
			Time:       0,
		},
		Cost:             7,
		PreviousPathNode: nil,
		NextQueueNode:    nil,
	}
	err = queue.Insert(&saturn) // Insert at second position
	assert.NoError(t, err)
	assert.Nil(t, earth.NextQueueNode)
	assert.NotNil(t, saturn.NextQueueNode)
	assert.NotNil(t, mars.NextQueueNode)
	assert.True(t, &mars == queue.firstItem, "Should have same addresses")
	assert.True(t, &saturn == queue.firstItem.NextQueueNode, "Should have same addresses")
	assert.True(t, &earth == queue.firstItem.NextQueueNode.NextQueueNode, "Should have same addresses")

	assert.Equal(t, []*DijkstraNode{&mars, &saturn, &earth}, queue.toSlice())
	assert.Len(t, queue.indexedNodes, 3)
	assert.True(t, queue.indexedNodes[earth.index] == &earth)
	assert.True(t, queue.indexedNodes[mars.index] == &mars)
	assert.True(t, queue.indexedNodes[saturn.index] == &saturn)

	err = queue.Insert(&mars) // Insert a node that is already there
	assert.Error(t, err)      // should raise an error

	sameIndex := DijkstraNode{
		index: mars.index,
		Cost:  125,
	}
	err = queue.Insert(&sameIndex) // Insert a node that is already there
	assert.Error(t, err)           // should raise an error

	err = queue.Remove(&saturn)
	assert.NoError(t, err)
	assert.Len(t, queue.indexedNodes, 2)
	assert.Equal(t, []*DijkstraNode{&mars, &earth}, queue.toSlice())

	err = queue.Remove(&mars)
	assert.NoError(t, err)
	assert.Len(t, queue.indexedNodes, 1)
	assert.Equal(t, []*DijkstraNode{&earth}, queue.toSlice())

	err = queue.Remove(&mars)
	assert.Error(t, err) // Removing a node that is not preset should raise an error

	err = queue.Remove(&earth)
	assert.NoError(t, err)
	assert.Len(t, queue.indexedNodes, 0)
	assert.Equal(t, []*DijkstraNode{}, queue.toSlice())
}

func BenchmarkInsert(b *testing.B) {
	nodes := make([]DijkstraNode, b.N)
	for i := 0; i < b.N; i++ {
		nodes[i].Cost = rand.Int()
		nodes[i].index = PlanetIndex{
			PlanetName: fake.City(),
			FuelLeft:   rand.Int(),
			Time:       rand.Int(),
		}
	}

	queue := DijkstraPriorityQueue{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = queue.Insert(&nodes[i])
	}
}

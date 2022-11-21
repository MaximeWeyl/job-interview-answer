package oddslib

// This file implements the priority queue of the dijkstra algorithm.

// ErrQueue is returned when a queue operation fails
type ErrQueue struct {
	message string
}

func (err ErrQueue) Error() string {
	return err.message
}
func (err ErrQueue) Unwrap() error {
	return ErrOddslib
}

// DijkstraPriorityQueue : Core queue of the dijkstra algorithm
// It is always sorted by the ascending cost
type DijkstraPriorityQueue struct {
	firstItem    *DijkstraNode
	indexedNodes map[PlanetIndex]*DijkstraNode
}

// Insert : insert a new element while making sure the queue remains sorted
func (queue *DijkstraPriorityQueue) Insert(nodeToInsert *DijkstraNode) (err error) {
	if nodeToInsert.NextQueueNode != nil {
		return ErrQueue{"the node to be inserted is already part of a queue"}
	}
	if queue.indexedNodes == nil {
		queue.indexedNodes = make(map[PlanetIndex]*DijkstraNode)
	}
	if _, ok := queue.indexedNodes[nodeToInsert.index]; ok {
		return ErrQueue{"a node with the same index already exists in the queue"}
	}

	// If the queue is empty, we insert as the first element
	if queue.firstItem == nil {
		queue.firstItem = nodeToInsert
		queue.indexedNodes[nodeToInsert.index] = nodeToInsert
		return
	}

	// We loop on the queue to find the right spot to insert
	var currentNode, previousNode *DijkstraNode = queue.firstItem, nil
	for ; currentNode != nil; previousNode, currentNode = currentNode, currentNode.NextQueueNode {
		// If we found the right spot
		if nodeToInsert.Cost < currentNode.Cost {
			switch currentNode {
			case queue.firstItem:
				nodeToInsert.NextQueueNode = currentNode
				queue.firstItem = nodeToInsert
			default:
				previousNode.NextQueueNode = nodeToInsert
				nodeToInsert.NextQueueNode = currentNode
			}

			queue.indexedNodes[nodeToInsert.index] = nodeToInsert
			return
		}
	}

	// If we did not find the right spot : we must insert at the end
	previousNode.NextQueueNode = nodeToInsert
	queue.indexedNodes[nodeToInsert.index] = nodeToInsert
	return
}

// Remove removes a node from the queue and from its index
// Other nodes links are updated too
func (queue *DijkstraPriorityQueue) Remove(node *DijkstraNode) (err error) {
	indexedNode, ok := queue.indexedNodes[node.index]
	if !ok || node != indexedNode {
		return ErrQueue{"tryied to remove a node from a queue that does not contains it"}
	}

	var currentNode, previousNode *DijkstraNode = queue.firstItem, nil
	for ; currentNode != node && currentNode != nil; previousNode, currentNode = currentNode, currentNode.NextQueueNode {
	}

	if currentNode == nil {
		panic("Fatal error : if this is triggered, there must be a bug inside the queue implementation")
	}

	if currentNode == queue.firstItem {
		queue.firstItem = currentNode.NextQueueNode
	} else {
		previousNode.NextQueueNode = currentNode.NextQueueNode
	}

	currentNode.NextQueueNode = nil
	delete(queue.indexedNodes, node.index)

	return
}

func (queue DijkstraPriorityQueue) toSlice() (out []*DijkstraNode) {
	out = []*DijkstraNode{}
	for current := queue.firstItem; current != nil; current = current.NextQueueNode {
		out = append(out, current)
	}
	return
}

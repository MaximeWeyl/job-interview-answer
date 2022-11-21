// Package oddslib : My implementation of the dijkstra algorithm
// See here for the inspiration : https://www.youtube.com/watch?v=GazC3A4OQTE&ab_channel=Computerphile
// This is a lazy version, which does not need to build the entire graph before running.
// Instead, it builds the nodes as it goes, which allows to use less memory and time
// The problem here is not to minimize the distance but the probability being caught
// The time dependent nature of the problem makes the weights on the edges depend on time
// (are the bounty hunters supposed to be here on this day or not ?).
// So the Nodes are indeed not just representing planets, but a (Planet, Fuel, Time) triplet.
// The graph is built accordingly to the problems specifications.
//
// An edge that leaves node A(PlanetA, FuelA, TimeA) to node B(PlanetB, FuelB, TimeB) must
// respect these constraints, or does not exist :
// - TimeB = TimeA + timeTravel(A->B) if TimeB <= countdown
// - FuelB = FuelA - timeTravel(A->B) if A != B, and timeTravel(A->B) <= FuelA
// - FuelB = FuelMAX if A == B
// This edge (and thus, NodeB) exist if and only if those constraints are respected.
//
// If node of arrival can be reached, we have found a path ! The dijkstra algorithm
// makes sure it is the cheapest in terms of probability being caught.
package oddslib

import (
	"fmt"
	"math"
)

// PlanetIndex stores the data that identifies a unique node into the algorithm
type PlanetIndex struct {
	PlanetName string
	FuelLeft   int
	Time       int
}

// DijkstraNode : A node for completing the Dijkstra algorithm
type DijkstraNode struct {
	index            PlanetIndex
	Cost             int
	PreviousPathNode *DijkstraNode
	NextQueueNode    *DijkstraNode
}

// Solve : Runs the actual path solving
func Solve(milleniumInfo MilleniumFalconInfo, galaxy Galaxy, empireInfo EmpireInfo) (oddsOfSuccess float64, solvedPath []PlanetIndex, err error) {
	departure := milleniumInfo.Departure
	arrival := milleniumInfo.Arrival

	// We check that it makes sense to even start a solving process
	// --> Falcon
	if !galaxy.HasPlanet(departure) {
		err = fmt.Errorf("The departure planet %s could not be found in the entire galaxy, where are you Solo ?\n%w", departure, ErrSolving)
		return
	}
	if !galaxy.HasPlanet(arrival) {
		err = fmt.Errorf("The arrival planet %s could not be found in the entire galaxy, let the empire waste its time...\n%w", arrival, ErrSolving)
		return
	}
	if milleniumInfo.Autonomy <= 0 {
		err = fmt.Errorf("Wait... Where do you think your going with this fuel tank of capacity %v\n%w", milleniumInfo.Autonomy, ErrSolving)
		return
	}
	// --> Empire
	if empireInfo.Countdown <= 0 {
		err = fmt.Errorf("Don't bother, it's gone already... (countdown is %v)\n%w", empireInfo.Countdown, ErrSolving)
		return
	}

	departureNode := DijkstraNode{
		index: PlanetIndex{
			PlanetName: departure,
			FuelLeft:   milleniumInfo.Autonomy,
			Time:       0,
		},
		Cost:             0,
		PreviousPathNode: nil,
	}

	bountyHunterChecker := NewBountyHunterChecker(empireInfo)

	queue := DijkstraPriorityQueue{}
	err = queue.Insert(&departureNode)
	if err != nil {
		fatal()
	}

	// Main loop : we always handle the top node of our priority queue, check the
	// youtube video to understand how it works (see top file comment)
	var currentNode *DijkstraNode
	for currentNode = queue.firstItem; currentNode != nil && currentNode.index.PlanetName != arrival; currentNode = queue.firstItem {
		// We look at all the places we could go to from this node
		for _, route := range galaxy.Routes[currentNode.index.PlanetName] {
			// Check the destination and update the queue accordingly
			destination := route.Destination
			time := currentNode.index.Time + route.TravelTime
			fuel := currentNode.index.FuelLeft - route.TravelTime
			handleDestination(&queue, currentNode, &bountyHunterChecker, destination, fuel, time, empireInfo.Countdown)
		}

		// We also consider staying here for the night
		// Note : what does a day even mean in the context of many different planets ?
		// :D, whatever...
		// Check the destination and update the queue accordingly
		destination := currentNode.index.PlanetName // We stay here
		time := currentNode.index.Time + 1          // For one night
		fuel := milleniumInfo.Autonomy              // And get a full refuel
		handleDestination(&queue, currentNode, &bountyHunterChecker, destination, fuel, time, empireInfo.Countdown)

		// We now have fully handled the current node
		// We remove it from the priority queue (see the video about the dijkstra algorithm)
		err = queue.Remove(currentNode)
		if err != nil {
			fatal()
		}
	}

	// If we cleared the queue, that means we never found the arrival
	// It is impossible to get to the arrival before the destruction
	switch {
	case currentNode == nil:
		oddsOfSuccess = 0
		solvedPath = []PlanetIndex{}
	case currentNode.index.PlanetName == arrival:
		oddsOfSuccess = math.Pow((9./10.), float64(currentNode.Cost))

		reversedPath := make([]PlanetIndex, 0)
		for node:= currentNode; node != nil ; node = node.PreviousPathNode {
			reversedPath = append(reversedPath, node.index)
		}
		solvedPath = make([]PlanetIndex, 0, len(reversedPath))
		for i := len(reversedPath)-1; i>=0 ; i-- {
			solvedPath = append(solvedPath, reversedPath[i])
		} 

	default:
		fatal()
	}

	return
}

func fatal() {
	panic("fatal error : this should happen due to how the algorithm works")
}

func handleDestination(
	queue *DijkstraPriorityQueue,
	currentNode *DijkstraNode,
	bountyHunterChecker *BountyHunterChecker,
	destination string,
	fuel, time, dayOfExplosion int) {

	var err error

	// If we do not have enough fuel to go there, we ignore this way
	if fuel < 0 {
		return
	}

	// If the death star has already with this way
	// We stop exploring this way : they're dead, what's the point...
	if time > dayOfExplosion {
		return
	}

	// We increment the cost if bounty hunters are found here on this particular time
	cost := currentNode.Cost
	if bountyHunterChecker.IsPlanetWatched(BountyHunterInfo{
		Planet: destination,
		Day:    time,
	}) {
		cost++
	}

	// If we got there : this node may be handled in the future, so we
	// insert it in out priority queue
	newNode := DijkstraNode{
		index: PlanetIndex{
			PlanetName: destination,
			FuelLeft:   fuel,
			Time:       time,
		},
		Cost:             cost,
		PreviousPathNode: currentNode,
	}

	if alreadyThereNode, ok := queue.indexedNodes[newNode.index]; ok {
		// If this node is already in the queue, we check if we improved
		// the cost, and update it only if we did
		if newNode.Cost < alreadyThereNode.Cost {
			// To update the node, we remove it and reinsert it, so the queue remains
			// sorted.
			err = queue.Remove(alreadyThereNode)
			if err != nil {
				fatal()
			}
			err = queue.Insert(&newNode)
			if err != nil {
				fatal()
			}
		}
	} else {
		// If this node is not in the queue, we insert it directly
		err = queue.Insert(&newNode)
		if err != nil {
			fatal()
		}
	}
}

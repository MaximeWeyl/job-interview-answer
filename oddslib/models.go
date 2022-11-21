package oddslib

// Galaxy : This is the set of all planets in the galaxy and how they are connected
// This is indexed by the origins. For an origin planets, we can get the list of all
// routes, sorted by the travel time.
type Galaxy struct {
	Routes map[string][]Route
}

// AddTwoWayRoute : Adds a two way route between two planets
func (galaxy *Galaxy) AddTwoWayRoute(route Route) {
	if galaxy.Routes == nil {
		galaxy.Routes = make(map[string][]Route)
	}

	// We run this one time for each two ways
	for i, planet := range []string{route.Origin, route.Destination} {
		planetRoutes, ok := galaxy.Routes[planet]
		if !ok {
			// If a planet is totally unknown in the index, we initialize it
			planetRoutes = make([]Route, 0)
			galaxy.Routes[planet] = planetRoutes
		}

		// We append this way
		if i > 0 {
			route.Origin, route.Destination = route.Destination, route.Origin
		}
		galaxy.Routes[planet] = append(planetRoutes, route)
	}
}

// HasPlanet : checks if the galaxy has a planet named like the provided string
func (galaxy Galaxy) HasPlanet(planetName string) bool {
	_, ok := galaxy.Routes[planetName]
	return ok
}


// BountyHunterChecker is used to check if a bounty hunter is expected 
// on a planet for a given day
type BountyHunterChecker struct {
	index map[BountyHunterInfo]struct{}
}

// IsPlanetWatched returns true if this planet is watched for this particular day
// If the planet is watched, the falcon could be caught !
func (checker BountyHunterChecker) IsPlanetWatched(index BountyHunterInfo) bool  {
	_, ok := checker.index[index]
	return ok
}

// NewBountyHunterChecker builds a new bounty hunter checker
func NewBountyHunterChecker(empireInfo EmpireInfo) BountyHunterChecker {
	checker := BountyHunterChecker{
		index: map[BountyHunterInfo]struct{}{},
	}

	for _, info := range empireInfo.BountyHunters {
		checker.index[info] = struct{}{}
	}

	return checker
}
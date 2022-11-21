package oddslib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGalaxy(t *testing.T) {
	galaxy := Galaxy{}

	route := Route{
		Origin:      "Earth",
		Destination: "Mars",
		TravelTime:  100,
	}
	galaxy.AddTwoWayRoute(route)

	earth, ok := galaxy.Routes["Earth"]
	assert.True(t, ok)
	assert.Len(t, earth, 1)
	assert.Equal(t, earth[0].Origin, "Earth")
	assert.Equal(t, earth[0].Destination, "Mars")
	assert.Equal(t, earth[0].TravelTime, 100)

	mars, ok := galaxy.Routes["Mars"]
	assert.True(t, ok)
	assert.Len(t, mars, 1)
	assert.Equal(t, mars[0].Origin, "Mars")
	assert.Equal(t, mars[0].Destination, "Earth")
	assert.Equal(t, mars[0].TravelTime, 100)
}

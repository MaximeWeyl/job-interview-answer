package oddslib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getSpecificationsResource(path string) string {
	return "../rebels-specifications/" + path
}

func TestLoadRoutes(t *testing.T) {
	dbPath := getSpecificationsResource("examples/example1/universe.db")

	infos := MilleniumFalconInfoWithPath{
		MilleniumFalconInfo: MilleniumFalconInfo{
			Autonomy:     12,
			Departure:    "PlanetA",
			Arrival:      "PlanetB",
		},
		RoutesDBPath: dbPath,
	}
	univese, err := infos.LoadGalaxy()
	assert.NoError(t, err)
	assert.NotEmpty(t, univese)
	assert.Len(t, univese.Routes, 4, "There must be 4 planets in this dataset")
}

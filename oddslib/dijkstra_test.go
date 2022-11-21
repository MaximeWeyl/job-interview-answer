package oddslib

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExamplesFromSpecifications(t *testing.T) {
	type Answer struct {
		Odds float64 `json:"odds"`
	}

	rootPath := getSpecificationsResource("examples/")
	for i := 1; i <= 4; i++ {
		t.Run(fmt.Sprintf("Example_{%v}", i), func(tt *testing.T) {
			path := rootPath + fmt.Sprintf("example%d/", i)

			var milleniumInfo MilleniumFalconInfoWithPath
			falconPath := path + "millenium-falcon.json"
			err := LoadJSONStructFromDisk(falconPath, &milleniumInfo)
			require.NoError(t, err)

			var empireInfo EmpireInfo
			err = LoadJSONStructFromDisk(path+"empire.json", &empireInfo)
			require.NoError(t, err)

			galaxy, err := milleniumInfo.LoadGalaxy(falconPath)
			require.NoError(t, err)
			require.NotEmpty(t, galaxy)

			var answer Answer
			err = LoadJSONStructFromDisk(path+"answer.json", &answer)
			require.NoError(t, err)

			resultOdds, _, err := Solve(milleniumInfo.MilleniumFalconInfo, galaxy, empireInfo)
			require.NoError(t, err)
			assert.Equal(t, answer.Odds, resultOdds)
		})
	}
}

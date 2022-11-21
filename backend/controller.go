package main

import (
	"errors"
	"net/http"

	"github.com/MaximeWeyl/misteryemployer-what-are-the-odds/oddslib"
	"github.com/gin-gonic/gin"
)

// GiveThemTheOddsInput is the json input that the route expects as a body
type GiveThemTheOddsInput struct {
	Empire          *oddslib.EmpireInfo          `json:"empire"`
	MilleniumFalcon *oddslib.MilleniumFalconInfo `json:"millenium-falcon"`
	Galaxy          *GalaxyInput                 `json:"galaxy"`
}

// GalaxyInput is the json representation of a galaxy
type GalaxyInput []oddslib.Route

// GiveThemTheOddsOutput is the json representation of the route output
type GiveThemTheOddsOutput struct {
	Odds float64               `json:"odds"`
	Path []oddslib.PlanetIndex `json:"path"`
}

// GiveThemTheOdds godoc
// @Summary Gives the odds of survival for a given problem
// @Description The input json must contain at least the 'empire' field
// The other fields are obsolete as this program comes with default
// millenium and galaxy data.
// @Param input body GiveThemTheOddsInput true "Only the 'empire' field is mandatory"
// @Success 200 {object} GiveThemTheOddsOutput "Get the odd and the solution"
// @Router /give-me-the-odds [post]
func GiveThemTheOdds(c *gin.Context) {
	// Parse the request body
	var input GiveThemTheOddsInput
	err := c.BindJSON(&input)
	if err != nil {
		handleError(c, err, http.StatusBadRequest)
	}

	// Retrives the mandatory empire, or fails
	if input.Empire == nil {
		err = errors.New("input json 'empire' is mandatory")
		handleError(c, err, http.StatusBadRequest)
	}
	empireInfo := *input.Empire

	// Use either the provided galaxy or the default
	var galaxy oddslib.Galaxy = DefaultGalaxy
	if input.Galaxy != nil {
		galaxy = oddslib.Galaxy{}
		for _, route := range *input.Galaxy {
			galaxy.AddTwoWayRoute(route)
		}
	}

	// Use either the provided falcon or the default
	milleniumInfo := DefaultMilleniumInfo
	if input.MilleniumFalcon != nil {
		milleniumInfo = *input.MilleniumFalcon
	}

	// Solves the problem
	odds, solvedPath, err := oddslib.Solve(milleniumInfo, galaxy, empireInfo)
	if err != nil {
		// Solving may end in an error
		handleError(c, err, http.StatusInternalServerError)
	}

	// Returns the final json to the client
	c.JSON(http.StatusOK, GiveThemTheOddsOutput{
		Odds: odds,
		Path: solvedPath,
	})
}

func handleError(c *gin.Context, err error, code int) {
	c.JSON(code, struct {
		Error string
	}{
		Error: err.Error(),
	})
}

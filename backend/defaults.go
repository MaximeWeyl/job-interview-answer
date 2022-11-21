package main

import (
	"fmt"

	"github.com/MaximeWeyl/misteryemployer-what-are-the-odds/oddslib"
	"github.com/spf13/cobra"
)

var (
	// DefaultMilleniumInfo is the default millenium falcon info used when
	// no info is given as input
	DefaultMilleniumInfo oddslib.MilleniumFalconInfo

	// DefaultGalaxy is the default galaxy info used when
	// no info is given as input
	DefaultGalaxy oddslib.Galaxy
)

func initFromCobra(cmd *cobra.Command, args []string) (err error) {
	milleniumPath := args[0]

	milleniumInfoWithPath := oddslib.MilleniumFalconInfoWithPath{}
	err = oddslib.LoadJSONStructFromDisk(milleniumPath, &milleniumInfoWithPath)
	if err != nil {
		err = fmt.Errorf("Error while trying to load default millenium info : %w", err)
		return
	}

	DefaultMilleniumInfo = milleniumInfoWithPath.MilleniumFalconInfo
	DefaultGalaxy, err = milleniumInfoWithPath.LoadGalaxy(milleniumPath)
	if err != nil {
		err = fmt.Errorf("Error while trying to load default galaxy : %w", err)
		return
	}

	return nil
}

package main

import (
	"fmt"
	"os"

	"github.com/MaximeWeyl/misteryemployer-what-are-the-odds/oddslib"
	"github.com/spf13/cobra"
)

// FlagShowHow is a CLI flag used to add the
// best path to the output
var FlagShowHow bool

func main() {
	cmd := cobra.Command{
		Use:   "r2d2 millenium_falcon_file empire_file",
		Short: "Finds a way that maximizes the odds of saving the galaxy, and displays the odds",
		Long: `
  .=.   --> Finds a way that maximizes the odds of saving the galaxy, and displays the odds
 '==c|  Biip bip biiiip bip
 [)-+|  Welcome to the R2-D2 CLI interface
 //'_|  Please provide the input data, and let
/]==;\  me help you save the galaxy.
		`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args[0], args[1])
		},
	}

	cmd.PersistentFlags().BoolVarP(&FlagShowHow, "show-how", "s", false, "verbose output")

	err := cmd.Execute()
	if err != nil {
		os.Exit(-1)
	}
}

func run(milleniumFalconFile, empireFile string) error {
	milleniumInfoWithPath := oddslib.MilleniumFalconInfoWithPath{}
	empireInfo := oddslib.EmpireInfo{}

	errFalcon := oddslib.LoadJSONStructFromDisk(milleniumFalconFile, &milleniumInfoWithPath)
	if errFalcon != nil {
		fmt.Println("There was an error while trying to parse the millenium falcon info")
		return errFalcon
	}

	errEmpire := oddslib.LoadJSONStructFromDisk(empireFile, &empireInfo)
	if errEmpire != nil {
		fmt.Println("There was an error while trying to parse the empire info")
		return errEmpire
	}

	galaxy, errGalaxy := milleniumInfoWithPath.LoadGalaxy(milleniumFalconFile)
	if errGalaxy != nil {
		fmt.Println("There was an error while trying import the galaxy")
		return errGalaxy
	}

	falconInfo := milleniumInfoWithPath.MilleniumFalconInfo
	result, solvedPath, errSolver := oddslib.Solve(falconInfo, galaxy, empireInfo)
	if errSolver != nil {
		fmt.Println("There was an error while trying to solve the problem")
		return errSolver
	}

	var comment string
	switch {
	case result == 0:
		comment = "They're screwed..."
	case result < 0.3:
		comment = "It's really difficult, send Luke"
	case result < 1:
		comment = "Be careful of bounty hunters"
	default:
		comment = "Easy mission"
	}

	fmt.Printf(
		"          ,-----.\n"+
			"         ,'_/_|_\\_`            BIP BIIIP BIP BIP BIIP BIP\n"+
			"        /<<::8[O]::>\\          BIP BIIIP BIP BBIP BIP BIP\n"+
			"       _|-----------|_         \n"+
			"   :::|  | ====-=- |  |:::     (\n"+
			"   :::|  | -=-==== |  |:::      Translates to : \n"+
			"   :::\\  | ::::|()||  /:::      Probability of success : %v %%\n"+
			"   ::::| | ....|()|| |::::      %v\n"+
			"       | |_________| |         )\n"+
			"       | |\\_______/| |\n"+
			"      /   \\ /   \\ /   \\   \n"+
			"      ``---' ``---' ``---'   \n", result*100, comment)

	if FlagShowHow {
		pathString := solvedPath[0].PlanetName
		arrowString := ""
		for i := 1; i < len(solvedPath); i++ {
			if solvedPath[i].PlanetName == solvedPath[i-1].PlanetName {
				arrowString = "R"
			} else {
				arrowString = fmt.Sprintf("%d", solvedPath[i].Time-solvedPath[i-1].Time)
			}

			pathString += fmt.Sprintf(
				" --%s--> %s",
				arrowString,
				solvedPath[i].PlanetName,
			)
		}
		fmt.Println("As requested, here is one optimal path")
		fmt.Println(pathString)
	}

	return nil
}

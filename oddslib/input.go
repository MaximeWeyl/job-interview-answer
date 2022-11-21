package oddslib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// EmpireInfo : Input info about the empire tactics
type EmpireInfo struct {
	Countdown     int                `json:"countdown"`
	BountyHunters []BountyHunterInfo `json:"bounty_hunters"`
}

// UnmarshalJSON parse a json and checks the mandatory fields
func (item *EmpireInfo) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		Countdown     *int                `json:"countdown"`
		BountyHunters *[]BountyHunterInfo `json:"bounty_hunters"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	}

	if required.Countdown == nil {
		err = fmt.Errorf("Field 'countdown' is required for Empire Info")
		return
	}
	item.Countdown = *required.Countdown

	if required.BountyHunters == nil {
		err = fmt.Errorf("Field 'bounty_hunters' is required for Empire Info")
		return
	}
	item.BountyHunters = *required.BountyHunters

	return
}

// BountyHunterInfo : Information about a specific event of a bounty hunter looking at a planet
//          |~
//          |.---.
//         .'_____`. /\
//         |~xxxxx~| ||
//         |_  #  _| ||
//    .------`-#-'-----.
//   (___|\_________/|_.`.
//    /  | _________ | | |
//   /   |/   _|_   \| | |
//  /   /X|  __|__  |/ `.|
// (  --< \\/    _\//|_ |`.
// `.    ~----.-~=====,:=======
//   ~-._____/___:__(``/| |
//     |    |      XX|~ | |
//      \__/======| /|  `.|
//      |_\|\    /|/_|    )
//      |_   \__/   _| .-'
//      | \ .'||`. / |(_|
//      |  ||.'`.||  |   )
//      |  `'|  |`'  |  /
//      |    |  |    |\/
type BountyHunterInfo struct {
	Planet string `json:"planet"`
	Day    int    `json:"day"`
}

// MilleniumFalconInfoWithPath : Information about our good old rusty starship !
//               c==o
//             _/____\_
//      _.,--'" ||^ || "`z._
//     /_/^ ___\||  || _/o\ "`-._
//   _/  ]. L_| || .||  \_/_  . _`--._
//  /_~7  _ . " ||. || /] \ ]. (_)  . "`--.
// |__7~.(_)_ []|+--+|/____T_____________L|
// |__|  _^(_) /^   __\____ _   _|
// |__| (_){_) J ]K{__ L___ _   _]
// |__| . _(_) \v     /__________|________
// l__l_ (_). []|+-+-<\^   L  . _   - ---L|
//  \__\    __. ||^l  \Y] /_]  (_) .  _,--'
//    \~_]  L_| || .\ .\\/~.    _,--'"
//     \_\ . __/||  |\  \`-+-<'"
//       "`---._|J__L|X o~~|[\\
//              \____/ \___|[//
type MilleniumFalconInfoWithPath struct {
	RoutesDBPath string `json:"routes_db"`
	MilleniumFalconInfo
}

// UnmarshalJSON parse a json and checks the mandatory fields
func (info *MilleniumFalconInfoWithPath) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		RoutesDBPath *string `json:"routes_db"`
		Autonomy     *int    `json:"autonomy"`
		Departure    *string `json:"departure"`
		Arrival      *string `json:"arrival"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	}

	if required.Arrival == nil {
		err = fmt.Errorf("Field 'arrival' is required for Millenium Falcon Info")
		return
	}
	info.Arrival = *required.Arrival

	if required.Departure == nil {
		err = fmt.Errorf("Field 'departure' is required for Millenium Falcon Info")
		return
	}
	info.Departure = *required.Departure

	if required.Autonomy == nil {
		err = fmt.Errorf("Field 'autonomy' is required for Millenium Falcon Info")
		return
	}
	info.Autonomy = *required.Autonomy

	if required.RoutesDBPath == nil {
		err = fmt.Errorf("Field 'routes_db_path' is required for Millenium Falcon Info")
		return
	}
	info.RoutesDBPath = *required.RoutesDBPath

	return
}

// LoadGalaxy : Uses the provided input to load the galaxy (or "universe")
func (info MilleniumFalconInfoWithPath) LoadGalaxy(basePaths ...string) (galaxy Galaxy, err error) {
	// Handles the fact that path could be relative or absolute
	basePath := ""
	if len(basePaths) > 1 {
		panic("Only one base path can be passed")
	}
	if len(basePaths) == 1 {
		basePath = filepath.Dir(basePaths[0])
	}
	filePath := filepath.Join(basePath, info.RoutesDBPath)

	// Handles the fact that the file may not exist
	if _, errStat := os.Stat(filePath); os.IsNotExist(errStat) {
		err = ErrFile{err: errStat}
		return
	}

	// Connect to the Sqlite DB
	db, err := gorm.Open(sqlite.Open(filePath), &gorm.Config{})
	if err != nil {
		err = DBError{err: err}
		return
	}

	// Loads the routes
	var routes []Route
	res := db.Find(&routes)
	if res.Error != nil {
		err = DBError{err: res.Error}
		return
	}

	// Builds the galaxy (wow)
	for _, route := range routes {
		galaxy.AddTwoWayRoute(route)
	}

	return
}

// MilleniumFalconInfo : This is the actual info about the falcon after
// removing the irrelevant DB path which the solver does not care about
type MilleniumFalconInfo struct {
	Autonomy  int    `json:"autonomy"`
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
}

// Route : This represents a route between two planets
// This is stored in a sqlite file
type Route struct {
	Origin      string
	Destination string
	TravelTime  int
}

// LoadJSONStructFromDisk loads a json file and assigns it to the target
func LoadJSONStructFromDisk(path string, target interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		cwd, _ := os.Getwd()
		return fmt.Errorf("%w, CWD was %s", err, cwd)
	}

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonBytes, target)
	if err != nil {
		return err
	}

	return nil
}

package database

import (
	"strings"
)

// need to build a search endpoitn so that we can serach for x and then return results
// in this example, we can search for bonds / coops and then return bonds/coops accoridngly
// this should be in the rpc package I guess, so don't do much until we take it there
func TestSearch(param string) {
	// param is the search string that you enter
	if strings.Contains(param, "bond") {
		// tkae bond actions
	} else if strings.Contains(param, "coop") {
		// do coop stuff
	}
}

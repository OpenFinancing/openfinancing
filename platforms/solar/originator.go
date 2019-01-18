package solar

import (
	"fmt"
	database "github.com/OpenFinancing/openfinancing/database"
	utils "github.com/OpenFinancing/openfinancing/utils"
)

// TODO: Describe Originators role and properties
// This function is returns a new entity with the bool "originator" as true
// TODO: Consider any other information needed for originators that should be added to the Users/Entities struct, or create a new struct altogether
func NewOriginator(uname string, pwd string, seedpwd string, Name string, Address string, Description string) (Entity, error) {
	return newEntity(uname, pwd, seedpwd, Name, Address, Description, "originator")
}

// Originate creates and saves a new origin contract
func (contractor *Entity) Originate(panelSize string, totalValue int, location string, years int, metadata string, recIndex int) (Project, error) {
	var pc Project
	var err error

	// for this, create a new  contract and store in the contracts db. Wea re sorting
	// by stage, so it shouldn't matter a whole lot
	indexCheck, err := RetrieveAllProjects()
	if err != nil {
		return pc, fmt.Errorf("Projects could not be retrieved!")
	}
	pc.Params.Index = len(indexCheck) + 1
	pc.Params.PanelSize = panelSize
	pc.Params.TotalValue = totalValue
	pc.Params.Location = location
	pc.Params.Years = years
	pc.Params.Metadata = metadata
	pc.Params.DateInitiated = utils.Timestamp()
	iRecipient, err := database.RetrieveRecipient(recIndex)
	if err != nil { // recipient does not exist
		return pc, err
	}
	pc.ProjectRecipient = iRecipient
	pc.Stage = 0 // 0 since we need to filter this out while retrieving the propsoed contracts
	pc.Originator = *contractor
	// instead of storing in this proposedcontracts slice, store it as a project, but not a contract and retrieve by stage
	err = pc.Save()
	// don't insert the project since the contractor's projects are not final
	return pc, err
}

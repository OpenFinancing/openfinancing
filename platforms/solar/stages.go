package solar

import (
	"fmt"
)

// TODO: get comments on the various stages involved here
var (
	PreOriginProject      = float64(0) // Stage 0: Originator approaches the recipient to originate an order. This is a project proposal.
	LegalContractStage    = 0.5        // Stage 0.5: Legal agreement (eg. MOU or letter of intent) between the originator and the recipient, out of blockchain. Can use a 2 of 2 multisig.
	OriginProject         = float64(1) // Stage 1: Originator/s proposes a contract on behalf of the recipient.
	OpenForMoneyStage     = 1.5        // Stage 1.5: The contract, even though not final, is now open to investors' money
	ProposedProject       = float64(2) // Stage 2: Contractors propose their contracts and investors can vote on them if they want to
	FinalizedProject      = float64(3) // Stage 3: Recipient chooses a particular contract for finalization. This can be arbitraty or following a specific tender process
	FundedProject         = float64(4) // Stage 4: Extend and Review the final legal contract, re-open for investment and finalize a particular contractor
	InstalledProjectStage = float64(5) // Stage 5: Installation of the panels / houses by the developer and contractor
	PowerGenerationStage  = float64(6) // Stage 6: Power generation and trigerring automatic payments, cover breach, etc.
	DebtPaidOffStage      = float64(7) // Stage 7: The stage at which the recipient pays back for his solar panels
)

// the following functions are helper functions to set the stage for a specific
// project
// we could also alternately define contract states and then read the state from
// our side and then compress this into a single function
func (a *Project) SetPreOriginProject() error {
	a.Stage = 0
	return a.Save()
}

func (a *Project) SetLegalContractStage() error {
	a.Stage = 0.5
	return a.Save()
}

func (a *Project) SetOriginProject() error {
	a.Stage = 1
	return a.Save()
}

func (a *Project) SetOpenForMoneyStage() error {
	a.Stage = 1.5
	return a.Save()
}

func (a *Project) SetProposedProject() error {
	a.Stage = 2
	return a.Save()
}

func (a *Project) SetFinalizedProject() error {
	a.Stage = 3
	return a.Save()
}

func (a *Project) SetFundedProject() error {
	a.Stage = 4
	return a.Save()
}

func (a *Project) SetInstalledProjectStage() error {
	a.Stage = 5
	return a.Save()
}

func (a *Project) SetPowerGenerationStage() error {
	a.Stage = 6
	return a.Save()
}

// MW: Consider the steps required for the promotion of the project to happen (eg. verification and validation)
func PromoteStage0To1Project(index int) error {
	// we need to upgrade the contract's whose index is contractIndex to stage 1
	projects, err := RetrieveProjectsAtStage(PreOriginProject)
	if err != nil {
		return err
	}
	fmt.Println("WE GET HERE", projects)
	for _, elem := range projects {
		fmt.Println("ELEM INDEX: ", elem.Params.Index)
		if elem.Params.Index == index {
			return elem.SetOriginProject() // upgrade stage of this project
		}
	}
	return fmt.Errorf("Project not found, erroring!")
}

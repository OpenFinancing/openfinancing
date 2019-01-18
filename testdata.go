package main

import (
	"fmt"
	"log"

	database "github.com/OpenFinancing/openfinancing/database"
	bonds "github.com/OpenFinancing/openfinancing/platforms/bonds"
	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
	utils "github.com/OpenFinancing/openfinancing/utils"
)

func InsertDummyData() error {
	var err error
	// populate database with dumym data
	var project1 solar.SolarParams
	var contract1 solar.Project
	var contract2 solar.Project
	var contract3 solar.Project
	var rec database.Recipient
	allRecs, err := database.RetrieveAllRecipients()
	if err != nil {
		log.Fatal(err)
	}
	if len(allRecs) == 0 {
		// there is no recipient right now, so create a dummy recipient
		var err error
		rec, err = database.NewRecipient("martin", "p", "x", "Martin")
		if err != nil {
			log.Fatal(err)
		}
	}

	var inv database.Investor
	allInvs, err := database.RetrieveAllInvestors()
	if err != nil {
		log.Fatal(err)
	}
	if len(allInvs) == 0 {
		var err error
		inv, err = database.NewInvestor("john", "p", "x", "John")
		if err != nil {
			log.Fatal(err)
		}
		err = inv.AddVotingBalance(100000)
		// this function saves as well, so there's no need to save again
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = bonds.NewBond("Dec 21 2021", "Member Rights Link", "Security Type 1", 5.4, "AAA", "Moody's Investments", "Wells Fargo",
		200000, "Opportunity Zone Construction", 200, "5% tax for 10 years", 1, "India Basin Project", "San Francisco", "India Basin is an upcoming creative project based in San Francisco that seeks to host innovators from all around the world")
	if err != nil {
		log.Fatal(err)
	}

	_, err = bonds.NewBond("Apr 2 2025", "Member Rights Link", "Security Type 2", 3.6, "AA", "Ant Financial", "People's Bank of China",
		50000, "Opportunity Zone Construction", 400, "No tax for 20 years", 1, "Shenzhen SEZ Development", "Shenzhen", "Shenzhen SEZ Development seeks to develop a SEZ in Shenzhen to foster creation of manufacturing jobs.")
	if err != nil {
		log.Fatal(err)
	}

	_, err = bonds.NewBond("Jul 9 2029", "Member Rights Link", "Security Type 3", 4.2, "BAA", "Softbank Corp.", "Bank of Japan",
		150000, "Opportunity Zone Construction", 100, "3% Tax for 5 Years", 1, "Osaka Development Project", "Osaka", "This Project seeks to develop cutting edge technologies in Osaka")
	if err != nil {
		log.Fatal(err)
	}
	// newParams(mdate string, mrights string, stype string, intrate float64, rating string, bIssuer string, uWriter string
	// unitCost float64, itype string, nUnits int, tax string
	coop, err := bonds.NewCoop("Dec 21 2021", "Member Rights Link", "Security Type 1", 5.4, "AAA", "Moody's Investments", "Wells Fargo",
		200000, "Coop Model", 4000, "India Basin Project", "San Francisco", "India Basin is an upcoming creative project based in San Francisco that seeks to host innovators from all around the world")
	if err != nil {
		log.Fatal(err)
	}

	_, err = bonds.NewCoop("Apr 2 2025", "Member Rights Link", "Security Type 2", 3.6, "AA", "Ant Financial", "People's Bank of China",
		50000, "Coop Model", 1000, "Shenzhen SEZ Development", "Shenzhen", "Shenzhen SEZ Development seeks to develop a SEZ in Shenzhen to foster creation of manufacturing jobs.")
	if err != nil {
		log.Fatal(err)
	}

	_, err = bonds.NewCoop("Jul 9 2029", "Member Rights Link", "Security Type 3", 4.2, "BAA", "Softbank Corp.", "Bank of Japan",
		150000, "Coop Model", 2000, "Osaka Development Project", "Osaka", "ODP seeks to develop cutting edge technologies in Osaka and invites investors all around the world to be a part of this new age")
	if err != nil {
		log.Fatal(err)
	}
	_, err = bonds.RetrieveCoop(coop.Params.Index)
	if err != nil {
		log.Fatal(err)
	}
	// NewOriginator(uname string, pwd string, Name string, Address string, Description string)
	newOriginator, err := solar.NewOriginator("john", "p", "x", "John Doe", "14 ABC Street London", "This is a sample originator")
	if err != nil {
		log.Fatal(err)
	}

	c1, err := solar.NewContractor("john", "p", "x", "John Doe", "14 ABC Street London", "This is a sample contractor")
	if err != nil {
		log.Println(err)
	}

	project1.Index = 1
	project1.PanelSize = "100 1000 sq.ft homes each with their own private spaces for luxury"
	project1.TotalValue = 14000
	project1.Location = "India Basin, San Francisco"
	project1.MoneyRaised = 0
	project1.Metadata = "India Basin is an upcoming creative project based in San Francisco that seeks to invite innovators from all around to participate"
	project1.InvestorAssetCode = ""
	project1.DebtAssetCode = ""
	project1.PaybackAssetCode = ""
	project1.DateInitiated = utils.Timestamp()
	project1.Years = 3
	contract1.Params = project1
	contract1.ProjectRecipient = rec
	contract1.Contractor = c1
	contract1.Originator = newOriginator
	contract1.Stage = 3
	err = contract1.Save()
	if err != nil {
		return fmt.Errorf("Error inserting project into db")
	}

	project1.Index = 2
	project1.PanelSize = "180 1200 sq.ft homes in a high rise building 0.1mi from Kendall Square"
	project1.TotalValue = 30000
	project1.Location = "Kendall Square, Boston"
	project1.MoneyRaised = 0
	project1.Metadata = "Kendall Square is set in the heart of Cambridge and is a popular startup IT hub"
	project1.InvestorAssetCode = ""
	project1.DebtAssetCode = ""
	project1.PaybackAssetCode = ""
	project1.DateInitiated = utils.Timestamp()
	project1.Years = 5
	contract2.ProjectRecipient = rec
	contract2.Params = project1
	contract2.Contractor = c1
	contract2.Originator = newOriginator
	contract2.Stage = 3
	err = contract2.Save()
	if err != nil {
		return fmt.Errorf("Error inserting project into db")
	}

	project1.Index = 3
	project1.PanelSize = "260 1500 sq.ft homes set in a medieval cathedral style construction"
	project1.TotalValue = 40000
	project1.Location = "Trafalgar Square, London"
	project1.MoneyRaised = 0
	project1.Metadata = "Trafalgar Square is set in the heart of London's financial district, with big banks all over"
	project1.InvestorAssetCode = ""
	project1.DebtAssetCode = ""
	project1.PaybackAssetCode = ""
	project1.DateInitiated = utils.Timestamp()
	project1.Years = 7
	contract3.ProjectRecipient = rec
	contract3.Params = project1
	contract3.Contractor = c1
	contract3.Originator = newOriginator
	contract3.Stage = 3
	err = contract3.Save()
	if err != nil {
		return fmt.Errorf("Error inserting project into db")
	}

	pc, err := newOriginator.Originate("100 16x24 panels on a solar rooftop", 14000, "Puerto Rico", 5, "ABC School in XYZ peninsula", 1) // 1 is the idnex for martin
	if err != nil {
		log.Fatal(err)
	}

	_, err = solar.RetrieveProject(pc.Params.Index)
	if err != nil {
		log.Fatal(err)
	}

	// Each contractor building off of this must reference the project index in their
	// proposed contract to enable searchability of the bucket. And each contractor
	// must build off of this in their proposed Contracts
	// Contractor stuff below, competing contractor details follow
	_, err = solar.NewContractor("sam", "p", "x", "Samuel Jackson", "14 ABC Street London", "This is a competing contractor")
	if err != nil {
		log.Fatal(err)
	}

	_, err = solar.NewOriginator("samuel", "p", "x", "Samuel L. Jackson", "ABC Street, London", "I am an originator")
	if err != nil {
		log.Fatal(err)
	}

	_, err = solar.RetrieveAllEntities("originator")
	if err != nil {
		log.Fatal(err)
	}
	_, err = solar.RetrieveAllEntities("contractor")
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

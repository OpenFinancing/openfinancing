// +build all travis

package solar

import (
	"log"
	"os"
	"testing"

	consts "github.com/OpenFinancing/openfinancing/consts"
	database "github.com/OpenFinancing/openfinancing/database"
)

// TODO: rewrite how this works and split between platforms and database
// go test --tags="all" -coverprofile=test.txt .
func TestDb(t *testing.T) {
	var err error
	os.Remove(os.Getenv("HOME") + "/.openfinancing/database/" + "/yol.db")
	err = os.MkdirAll(os.Getenv("HOME")+"/.openfinancing/database", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	oldDbDir := consts.DbDir
	consts.DbDir = "blah" // set to a false db so that we can test errors arising from OpenDB()
	x1, err := newEntity("OrigTest", "pwd", "blah", "NameOrigTest", "123 ABC Street", "OrigDescription", "originator")
	if err == nil {
		t.Fatalf("Able to create entity with invalid db, quitting!")
	}

	_, err = x1.Propose("100 16x32 panels", 28000, "Puerto Rico", 6, "LEED+ Gold rated panels and this is random data out of nowhere and we supply our own devs and provide insurance guarantee as well. Dual audit maintenance upto 1 year. Returns capped as per defaults", 1, 1)
	// 1 for retrieving martin as the recipient and 1 is the project Index
	if err == nil {
		t.Fatalf("Able to propose contract with invalid db, quitting!")
	}
	_, err = x1.Originate("100 16x24 panels on a solar rooftop", 14000, "Puerto Rico", 5, "ABC School in XYZ peninsula", 1) // 1 is the index for martin
	if err == nil {
		t.Fatal("Able to originate contract with invalid db, quitting!")
	}
	err = PromoteStage0To1Project(1)
	if err == nil {
		t.Fatalf("Able to promote contract even with invalid db, quitting!")
	}
	var y1 Project
	err = y1.Save()
	if err == nil {
		t.Fatalf("Able to save file even though no db is present, quitting!")
	}
	_, err = RetrieveProject(1)
	if err == nil {
		t.Fatalf("Able to retrieve project with invalid db, quitting!")
	}
	_, err = RetrieveAllProjects()
	if err == nil {
		t.Fatalf("Able to retrieve projects with invalid db, quitting!")
	}
	_, err = RetrieveProjectsAtStage(1)
	if err == nil {
		t.Fatalf("Able to retrieve stage projects with invalid db, quitting!")
	}
	_, err = RetrieveContractorProjects(1, 1)
	if err == nil {
		t.Fatalf("Able to retrieve contractor projects with invalid db, quitting!")
	}
	_, err = RetrieveOriginatorProjects(1, 1)
	if err == nil {
		t.Fatalf("Able to retrieve originated projects with invalid db, quitting!")
	}
	_, err = RetrieveRecipientProjects(1, 1)
	if err == nil {
		t.Fatalf("Able to retrieve recipient projects with invalid db, quitting!")
	}
	var xy1 database.Investor
	err = VoteTowardsProposedProject(xy1.U.Index, 1, 1)
	if err == nil {
		t.Fatalf("Can vote towards a non existent proposed project, quitting!")
	}
	var xyz1 Entity
	err = xyz1.Save()
	if err == nil {
		t.Fatalf("Can save entity which doesn't exist?")
	}
	_, err = RetrieveAllEntities("guarantor")
	if err == nil {
		t.Fatalf("Can retreive contract entities from invalid db, quitting!")
	}
	_, err = RetrieveEntity(1)
	if err == nil {
		t.Fatalf("Can retrieve entity in invalid db, quitting!")
	}
	_, err = SearchForEntity("invalid", "invalid")
	if err == nil {
		t.Fatalf("Can search for enttiy in invalid db, quitting!")
	}
	consts.DbDir = oldDbDir
	db, err := database.OpenDB()
	if err != nil {
		t.Fatal(err)
	}
	db.Close() // close immmediately after check
	var project1 SolarParams
	var dummy Project
	// investors entity testing over, test recipients below in the same way
	// now we repeat the same tests for all other entities
	// connections and the other for non RPC connections
	// now we repeat the same tests for all other entities
	// connections and the other for non RPC connections
	inv, err := database.NewInvestor("inv1", "blah", "blah", "cool")
	if err != nil {
		t.Fatal(err)
	}
	err = inv.Save()
	if err != nil {
		t.Fatal(err)
	}
	recp, err := database.NewRecipient("user1", "blah", "blah", "cool")
	if err != nil {
		t.Fatal(err)
	}
	err = recp.Save()
	if err != nil {
		t.Fatal(err)
	}

	newCE2, err := NewOriginator("OrigTest2", "pwd", "blah", "NameOrigTest2", "123 ABC Street", "OrigDescription2")
	if err != nil {
		t.Fatal(err)
	}
	// tests for contractors
	newCE, err := NewContractor("ConTest", "pwd", "blah", "NameConTest", "123 ABC Street", "ConDescription") // use and test this as well
	if err != nil {
		t.Fatal(err)
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
	project1.DateInitiated = ""
	project1.Years = 3
	dummy.Params = project1
	dummy.ProjectRecipient = recp
	dummy.Contractor = newCE
	dummy.Originator = newCE2
	dummy.Stage = 3
	err = dummy.Save()
	if err != nil {
		t.Errorf("Inserting an project into the database failed")
		// shouldn't really fatal here, but this is in main, so we can't return
	}
	project, err := RetrieveProject(dummy.Params.Index)
	if err != nil {
		log.Println(err)
		t.Errorf("Retrieving project from the database failed")
		// again, shouldn't really fat a here, but we're in main
	}
	klx1, err := RetrieveProject(1000)
	if klx1.Params.Index != 0 {
		t.Fatalf("REtrieved project which does not exist, quitting!")
	}
	if project.Params.Index != dummy.Params.Index {
		t.Fatalf("Indices don't match, quitting!")
	}
	dummy.Params.Index = 2 // change index and try inserting another project
	err = dummy.Save()
	if err != nil {
		log.Println(err)
		t.Errorf("Inserting an project into the database failed")
		// shouldn't really fatal here, but this is in main, so we can't return
	}
	projects, err := RetrieveAllProjects()
	if err != nil {
		log.Println("Retrieve all error: ", err)
		t.Errorf("Failed in retrieving all projects")
	}
	if projects[0].Params.Index != 1 {
		t.Fatalf("Index of first element doesn't match, quitting!")
	}
	oProjects, err := RetrieveProjectsAtStage(OriginProject)
	if err != nil {
		log.Println("Retrieve all error: ", err)
		t.Errorf("Failed in retrieving all originated projects")
	}
	if len(oProjects) != 0 {
		log.Println("OPROJECTS: ", len(oProjects))
		t.Fatalf("Originated projects present!")
	}
	err = database.DeleteKeyFromBucket(dummy.Params.Index, database.ProjectsBucket)
	if err != nil {
		log.Println(err)
		t.Errorf("Deleting an order from the db failed")
	}
	// err = DeleteProject(1000) this would work becuase the key will not be found and hence would not return an error
	pp1, err := RetrieveProject(dummy.Params.Index)
	if err == nil && pp1.Params.Index != 0 {
		log.Println(err)
		// this should fail because we're trying to read an empty key value pair
		t.Errorf("Found deleted entry, quitting!")
	}

	inv, err = database.NewInvestor("inv2", "b921f75437050f0f7d2caba6303d165309614d524e3d7e6bccf313f39d113468d30e1e2ac01f91f6c9b66c083d393f49b3177345311849edb026bb86ee624be0", "blah", "cool")
	if err != nil {
		t.Fatal(err)
	}
	err = inv.Save()
	if err != nil {
		t.Fatal(err)
	}
	recp, err = database.NewRecipient("user1", "b921f75437050f0f7d2caba6303d165309614d524e3d7e6bccf313f39d113468d30e1e2ac01f91f6c9b66c083d393f49b3177345311849edb026bb86ee624be0", "blah", "cool")
	if err != nil {
		t.Fatal(err)
	}
	err = recp.Save()
	if err != nil {
		t.Fatal(err)
	}
	// tests for originators
	newCE, err = newEntity("OrigTest", "pwd", "blah", "NameOrigTest", "123 ABC Street", "OrigDescription", "originator")
	if err != nil {
		t.Fatal(err)
	}
	allOrigs, err := RetrieveAllEntities("originator")
	if err != nil {
		t.Fatal(err)
	}
	acz1, err := RetrieveAllEntities("random")
	if len(acz1) != 0 {
		log.Println(acz1)
		t.Fatalf("Category which does not exist exists?")
	}
	if len(allOrigs) != 2 || allOrigs[0].U.Name != "NameOrigTest2" {
		t.Fatal("Names don't match, quitting!")
	}
	_, err = RetrieveAllEntities("guarantor")
	if err != nil {
		t.Fatal(err)
	}
	allConts, err := RetrieveAllEntities("contractor")
	if err != nil {
		t.Fatal(err)
	}
	if len(allConts) != 1 || allConts[0].U.Name != "NameConTest" {
		t.Fatal("Names don't match, quitting!")
	}
	_, err = newCE.Propose("100 16x32 panels", 28000, "Puerto Rico", 6, "LEED+ Gold rated panels and this is random data out of nowhere and we supply our own devs and provide insurance guarantee as well. Dual audit maintenance upto 1 year. Returns capped as per defaults", recp.U.Index, 1)
	// 1 for retrieving martin as the recipient and 1 is the project Index
	if err != nil {
		t.Fatal(err)
	}
	_, err = newCE.Propose("100 16x32 panels", 28000, "Puerto Rico", 6, "LEED+ Gold rated panels and this is random data out of nowhere and we supply our own devs and provide insurance guarantee as well. Dual audit maintenance upto 1 year. Returns capped as per defaults", 1000, 1)
	// 1 for retrieving martin as the recipient and 1 is the project Index
	if err == nil {
		t.Fatal("Able to retrieve non existent recipient, quitting!")
	}
	rOx, err := RetrieveProject(2)
	if err != nil {
		t.Fatal(err)
	}
	rOx.ProjectRecipient = recp
	err = rOx.Save()
	if err != nil {
		t.Fatal(err)
	}

	allPCs, err := RetrieveRecipientProjects(ProposedProject, 6)
	if err != nil {
		t.Fatal(err)
	}
	if len(allPCs) != 1 { // add check for stuff here
		log.Println("LEN all proposed projects", len(allPCs))
	}
	rPC, err := findInKey(2, allPCs)
	if err != nil {
		t.Fatal(err)
	}
	if rPC.Params.Index != rOx.Params.Index {
		t.Fatal("Indices don't match")
	}

	// now come the failure cases which should fail and we shall catch the case when they don't
	allPCs, err = RetrieveContractorProjects(ProposedProject, 2)
	if len(allPCs) != 0 {
		log.Println("LEBNGRG: ", len(allPCs))
		t.Fatalf("Retrieving a missing contract succeeds, quitting!")
	}
	rPC, err = findInKey(2, allPCs)
	if err == nil {
		t.Fatalf("Entity which should be missing exists!")
	}

	// newCE, err = newEntity("ConTest", "pwd", "NameConTest", "123 ABC Street", "ConDescription", "contractor")
	rCE, err := SearchForEntity("ConTest", "9f88a8d40b90616715f868ed195d24e5df994f56bce34eddb022c213484eb0f220d8907e4ecd8f64ddd364cb30bb5758b32ee26541f340b930f7e5bf756907a4")
	if err != nil {
		t.Fatal(err)
	}
	if rCE.Contractor != true {
		log.Println("The role: ", rCE.Contractor, rCE)
		t.Fatal("Roles don't match!")
	}
	trC1, err := RetrieveEntity(7)
	if err != nil || trC1.U.Index == 0 {
		t.Fatal("Project Entity lookup failed")
	}
	tmpx1, err := newCE2.Originate("100 16x24 panels on a solar rooftop", 14000, "Puerto Rico", 5, "ABC School in XYZ peninsula", recp.U.Index) // 1 is the index for martin
	if err != nil {
		t.Fatal(err)
	}
	_, err = newCE2.Originate("100 16x24 panels on a solar rooftop", 14000, "Puerto Rico", 5, "ABC School in XYZ peninsula", 1000) // 1 is the index for martin
	if err == nil {
		t.Fatalf("Not quitting for invalid recipient index")
	}
	tmpx1.Stage = 6
	err = tmpx1.Save()
	if err != nil {
		t.Fatal(err)
	}
	allOOs, err := RetrieveProjectsAtStage(6) // this checks for stage 1 and not zero like the thing installed above
	if err != nil {
		t.Fatal(err)
	}
	if len(allOOs) != 1 {
		log.Println("Length of all Stage 6 Projects: ", len(allOOs))
		t.Fatalf("Length of all stage 6 projects doesn't match")
	}
	indexCheck, err := RetrieveAllProjects()
	if err != nil {
		t.Fatalf("Projects could not be retrieved!")
	}
	project1.Index = len(indexCheck) + 1
	nOP, err := newOriginProject(project1, newCE2)
	if err != nil {
		t.Fatal(err)
	}
	nOP.ProjectRecipient = recp
	err = nOP.Save()
	if err != nil {
		t.Fatal(err)
	}
	err = nOP.SetPreOriginProject()
	if err != nil {
		t.Fatal(err)
	}

	err = nOP.SetLegalContractStage()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 0.5 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	err = nOP.SetOriginProject()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 1 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	err = nOP.SetOpenForMoneyStage()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 1.5 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	err = nOP.SetProposedProject()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 2 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	err = nOP.SetFinalizedProject()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 3 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	err = nOP.SetFundedProject()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 4 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	err = nOP.SetInstalledProjectStage()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 5 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	err = nOP.SetPowerGenerationStage()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 6 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	// cycle back to stage 0 and try using the other function to modify the stage
	err = nOP.SetPreOriginProject()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 0 {
		t.Fatalf("Stage doesn't match, quitting!")
	}
	err = PromoteStage0To1Project(1000)
	if err == nil {
		t.Fatalf("Able to promote project which doesn't exist!")
	}
	err = PromoteStage0To1Project(nOP.Params.Index)
	if err != nil {
		t.Fatal(err)
	}
	allO, err := RetrieveOriginatorProjects(OriginProject, newCE2.U.Index)
	if err != nil {
		t.Fatal(err)
	}
	if len(allO) != 1 {
		t.Fatalf("Multiple originated orders when there should be only one")
	}
	err = inv.AddVotingBalance(1000)
	if err != nil {
		t.Fatal(err)
	}
	err = nOP.SetProposedProject()
	if err != nil {
		t.Fatal(err)
	}
	err = VoteTowardsProposedProject(inv.U.Index, 100, 2)
	if err != nil {
		t.Fatal(err)
	}
	recp.ReceivedSolarProjects = append(recp.ReceivedSolarProjects, project1.DebtAssetCode)
	// the above thing is to test the function itself and not the functionality since
	// DebtAssetCode for project1 should be empty
	err = recp.Save()
	if err != nil {
		t.Fatal(err)
	}
	err = RecipientAuthorize(nOP.Params.Index, recp) // again a placeholder for testing this function, not the flow itself
	if err != nil {
		t.Fatal(err)
	}
	chk := nOP.CalculatePayback("100")
	if chk != "0.257143" {
		t.Fatalf("Balance doesn't match , quitting!")
	}
	var arr []Project
	x, err := SelectContractByPrice(arr)
	if err == nil {
		t.Fatalf("Empty array returns choice")
	}
	y, err := SelectContractByTime(arr)
	if err == nil {
		t.Fatalf("Empty array returns choice")
	}
	arr = append(arr, nOP)
	x, err = SelectContractByPrice(arr)
	if err != nil {
		t.Fatal(err)
	}
	sc1, err := RetrieveAllProjects()
	if err != nil {
		t.Fatal(err)
	}
	/*
		sc1[0]: YEARS: 3, PRICE: 14000
		sc1[1]: YEARS: 6, PRICE: 28000
		sc1[2]: YEARS: 5, PRICE: 14000
		sc1[3]: YEARS: 3, PRICE: 14000
	*/
	var arrx []Project // (6, 28000), (3, 14000)
	arrx = append(arr, sc1[1])
	arrx = append(arr, sc1[0])
	_, err = SelectContractByTime(arrx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = SelectContractByPrice(arrx)
	if err != nil {
		t.Fatal(err)
	}
	if x.Params.Index != nOP.Params.Index {
		t.Fatalf("Indices don't match, quitting!")
	}
	y, err = SelectContractByTime(arr)
	if err != nil {
		t.Fatal(err)
	}
	if y.Params.Index != nOP.Params.Index {
		t.Fatalf("Indices don't match, quitting!")
	}
	_, err = newEntity("OrigTest", "pwd", "blah", "NameOrigTest", "123 ABC Street", "OrigDescription", "developer")
	if err != nil {
		t.Fatal(err)
	}
	_, err = newEntity("OrigTest", "pwd", "blah", "NameOrigTest", "123 ABC Street", "OrigDescription", "guarantor")
	if err != nil {
		t.Fatal(err)
	}
	_, err = newEntity("OrigTest", "pwd", "blah", "NameOrigTest", "123 ABC Street", "OrigDescription", "invalid")
	if err == nil {
		t.Fatalf("Not able to catch invalid contractor error, quitting!")
	}
	_, err = RetrieveAllEntities("developer")
	if err != nil {
		t.Fatal(err)
	}
	_, err = RetrieveAllEntities("gurantor")
	if err != nil {
		t.Fatal(err)
	}
	_, err = SearchForEntity("name", "invalid")
	if err == nil {
		t.Fatalf("Not able to catch invalid entity error, quitting!")
	}
	// func Payback(a *database.Recipient, project Project, assetName string, issuerPubkey string, amount string, seed string) error
	err = Payback(1, 1, "", "", "", "", "")
	if err == nil {
		t.Fatal("Invalid params not caught, exiting!")
	}
	_, err = newEntity("x", "x", "x", "x", "123 ABC Street", "x", "random")
	if err == nil {
		t.Fatalf("not able to catch invalid entity error")
	}
	var recpx database.Recipient
	recpx.ReceivedSolarProjects = append(recpx.ReceivedSolarProjects, dummy.Params.DebtAssetCode)
	err = dummy.updateRecipient(recpx)
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(os.Getenv("HOME") + "/.openfinancing/database/" + "/yol.db")
}

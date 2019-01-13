package main

import (
	"fmt"
	"log"

	database "github.com/OpenFinancing/openfinancing/database"
	platform "github.com/OpenFinancing/openfinancing/platforms"
	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
	scan "github.com/OpenFinancing/openfinancing/scan"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
)

func ValidateInputs() {
	if (opts.RecYears != 0) && !(opts.RecYears == 3 || opts.RecYears == 5 || opts.RecYears == 7) {
		log.Fatal(fmt.Errorf("Number of years not supported"))
	}
}

func StartPlatform() (string, string, error) {
	var publicKey string
	var seed string
	ValidateInputs()
	database.CreateHomeDir()
	allContracts, err := solar.RetrieveAllProjects()
	if err != nil {
		log.Println("Error retrieving all projects from the database")
		return publicKey, seed, err
	}

	if len(allContracts) == 0 {
		log.Println("Populating database with test values")
		err = InsertDummyData()
		if err != nil {
			return publicKey, seed, err
		}
	}
	publicKey, seed, err = platform.InitializePlatform()
	return publicKey, seed, err
}

func NewUserPrompt() (string, string, string, string, error) {
	realName, err := scan.ScanForString()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return "", "", "", "", err
	}
	fmt.Printf("%s: ", "ENTER YOUR USERNAME")
	loginUserName, err := scan.ScanForString()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return "", "", "", "", err
	}

	fmt.Printf("%s: ", "ENTER DESIRED PASSWORD, YOU WILL NOT BE ASKED TO CONFIRM THIS")
	loginPassword, err := scan.ScanForPassword()
	if err != nil {
		fmt.Println("Couldn't read password")
		return "", "", "", "", err
	}
	fmt.Printf("%s: ", "ENTER SEED PASSWORD, YOU WILL NOT BE ASKED TO CONFIRM THIS")
	seedPassword, err := scan.ScanForPassword()
	return realName, loginUserName, loginPassword, seedPassword, err
}

func NewInvestorPrompt() error {
	log.Println("You have chosen to create a new investor account, welcome")
	loginUserName, loginPassword, realName, seedpwd, err := NewUserPrompt()
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = database.NewInvestor(loginUserName, loginPassword, seedpwd, realName)
	if err != nil {
		log.Println("FAILED TO SETUP ACCOUNT, TRY AGAIN")
		return err
	}
	return err
}

func NewRecipientPrompt() error {
	log.Println("You have chosen to create a new recipient account, welcome")
	loginUserName, loginPassword, realName, seedpwd, err := NewUserPrompt()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = database.NewRecipient(loginUserName, loginPassword, seedpwd, realName)
	if err != nil {
		log.Println("FAILED TO SETUP ACCOUNT, TRY AGAIN")
		return err
	}
	return err
}

func LoginPrompt() (database.Investor, database.Recipient, solar.Entity, bool, bool, error) {
	rbool := false
	cbool := false
	var investor database.Investor
	var recipient database.Recipient
	var contractor solar.Entity
	fmt.Println("---------SELECT YOUR ROLE---------")
	fmt.Println(" i. INVESTOR")
	fmt.Println(" r. RECIPIENT")
	fmt.Println(" c. CONTRACTOR")
	optS, err := scan.ScanForString()
	if err != nil {
		log.Println("Failed to read user input")
		return investor, recipient, contractor, rbool, cbool, err
	}
	if optS == "I" || optS == "i" {
		fmt.Println("WELCOME BACK INVESTOR")
	} else if optS == "R" || optS == "r" {
		fmt.Println("WELCOME BACK RECIPIENT")
		rbool = true
	} else if optS == "C" || optS == "c" {
		cbool = true
		fmt.Println("WELCOME BACK CONTRACTOR")
	} else {
		log.Println("INVALID INPUT, EXITING!")
		return investor, recipient, contractor, rbool, cbool, fmt.Errorf("INVALID INPUT, EXITING!")
	}
	// ask for username and password combo here
	fmt.Printf("%s: ", "ENTER YOUR USERNAME")
	loginUserName, err := scan.ScanForString()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return investor, recipient, contractor, rbool, cbool, err
	}

	fmt.Printf("%s: ", "ENTER YOUR PASSWORD: ")
	loginPassword, err := scan.ScanForPassword()
	if err != nil {
		fmt.Println("Couldn't read password")
		return investor, recipient, contractor, rbool, cbool, err
	}
	user, err := database.ValidateUser(loginUserName, loginPassword)
	if err != nil {
		fmt.Println("Couldn't read password")
		return investor, recipient, contractor, rbool, cbool, err
	}
	log.Println("WATCH USER INDEX: ", user.Index)
	if rbool {
		recipient, err = database.RetrieveRecipient(user.Index)
		if err != nil {
			return investor, recipient, contractor, rbool, cbool, err
		}
	} else if cbool {
		contractor, err = solar.RetrieveEntity(user.Index)
		if err != nil {
			return investor, recipient, contractor, rbool, cbool, err
		}
	} else {
		investor, err = database.RetrieveInvestor(user.Index)
		if err != nil {
			return investor, recipient, contractor, rbool, cbool, err
		}
	}
	return investor, recipient, contractor, rbool, cbool, nil
}

func OriginContractPrompt(contractor *solar.Entity) error {
	fmt.Println("YOU HAVE DECIDED TO PROPOSE A NEW CONTRACT")
	fmt.Println("ENTER THE PANEL SIZE")
	panelSize, err := scan.ScanForString()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE COST OF PROJECT")
	totalValue, err := scan.ScanForInt()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE LOCATION OF PROJECT")
	location, err := scan.ScanForString()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE ESTIMATED NUMBER OF YEARS FOR COMPLETION")
	years, err := scan.ScanForInt()
	if err != nil {
		return err
	}
	fmt.Println("ENTER METADATA REGARDING THE PROJECT")
	metadata, err := scan.ScanForString()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE RECIPIENT'S USER ID")
	recIndex, err := scan.ScanForInt()
	if err != nil {
		return err
	}
	originContract, err := contractor.OriginContract(panelSize, totalValue, location, years, metadata, recIndex)
	if err != nil {
		return err
	}
	// project insertion is done by the  above function, so we needn't call the database to do it again for us
	PrintProject(originContract)
	return nil
}

func ProposeContractPrompt(contractor *solar.Entity) error {
	fmt.Println("YOU HAVE DECIDED TO PROPOSE A NEW CONTRACT")
	fmt.Println("ENTER THE PROJECT INDEX")
	contractIndex, err := scan.ScanForInt()
	if err != nil {
		return err
	}
	// we need to check if this contract index exists and retrieve
	rContract, err := solar.RetrieveProject(contractIndex)
	if err != nil {
		return err
	}
	// TODO: PRint CONTRACTS HERE
	log.Println("YOUR CONTRACT IS: ")
	PrintProject(rContract)
	if rContract.Params.Index == 0 || rContract.Stage != 1 {
		// prevent people form porposing contracts for non originated contracts
		return fmt.Errorf("Invalid contract index")
	}
	panelSize := rContract.Params.PanelSize
	location := rContract.Params.Location
	fmt.Println("ENTER THE COST OF PROJECT")
	totalValue, err := scan.ScanForInt()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE ESTIMATED NUMBER OF YEARS FOR COMPLETION")
	years, err := scan.ScanForInt()
	if err != nil {
		return err
	}
	fmt.Println("ENTER METADATA REGARDING THE PROJECT")
	metadata, err := scan.ScanForString()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE RECIPIENT'S USER ID")
	recIndex, err := scan.ScanForInt()
	if err != nil {
		return err
	}
	originContract, err := contractor.ProposeContract(panelSize, totalValue, location, years, metadata, recIndex, contractIndex)
	if err != nil {
		return err
	}
	// project insertion is done by the  above function, so we needn't call the database to do it again for us
	PrintProject(originContract)
	return nil
}

func Stage3ProjectsDisplayPrompt() {
	fmt.Println("------------LIST OF ALL AVAILABLE PROJECTS------------")
	allProjects, err := solar.RetrieveProjects(solar.FinalizedProject)
	if err != nil {
		log.Println("Error retrieving all projects from the database")
	} else {
		PrintProjects(allProjects)
	}
}

func DisplayOriginProjects() {
	fmt.Println("PRINTING ALL ORIGINATED PROJECTS: ")
	x, err := solar.RetrieveProjects(solar.OriginProject)
	if err != nil {
		log.Println(err)
	} else {
		PrintProjects(x)
	}
}

func ExitPrompt() {
	// check whether he wants to go back to the display all screen again
	fmt.Println("DO YOU REALLY WANT TO EXIT? (PRESS Y TO CONFIRM)")
	exitOpt, err := scan.ScanForString()
	if err != nil {
		log.Println(err)
	}
	if exitOpt == "Y" || exitOpt == "y" {
		fmt.Println("YOU HAVE DECIDED TO EXIT")
		log.Fatal("")
	}
}

func BalanceDisplayPrompt(publicKey string) {
	balances, err := xlm.GetAllBalances(publicKey)
	if err != nil {
		log.Println(err)
	} else {
		PrintBalances(balances)
	}
}

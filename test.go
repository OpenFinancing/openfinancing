package main

// test.go runs the PoC stellar implementation calling various functions
import (
	"fmt"
	"log"
	"os"

	assets "github.com/OpenFinancing/openfinancing/assets"
	consts "github.com/OpenFinancing/openfinancing/consts"
	ipfs "github.com/OpenFinancing/openfinancing/ipfs"
	platform "github.com/OpenFinancing/openfinancing/platforms"
	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
	rpc "github.com/OpenFinancing/openfinancing/rpc"
	scan "github.com/OpenFinancing/openfinancing/scan"
	stablecoin "github.com/OpenFinancing/openfinancing/stablecoin"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
	flags "github.com/jessevdk/go-flags"
)

// test.go drives the CLI interface. Will be removed once we have a funcitoning frontend
// that supplements the backend in an effective way.
var opts struct {
	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	// TODO: define default values for each and then use them if not passed
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	// InvAmount int    `short:"i" description:"Desired investment" required:"true"`
	InvAmount int `short:"i" description:"Desired investment"`
	// RecYears  int    `short:"r" description:"Number of years the recipient wants to repay in. Can be 3, 5 or 7 years." required:"true"`
	RecYears int    `short:"r" description:"Number of years the recipient wants to repay in. Can be 3, 5 or 7 years."`
	Port     string `short:"p" description:"The port on which the server runs on"`
}

func main() {
	var err error
	_, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	go rpc.StartServer("8080") // run as go routine for now
	var investorSeed string
	var recipientSeed string
	/*
		fmt.Println("WHICH PLATFORM WOULD YOU IKE TO ENTER INTO?")
		fmt.Println("1. Platform of Contracts")
		fmt.Println("1. Platform of Platforms")
		platformArg, err := scan.ScanForInt()
		if err != nil {
			log.Fatal(err)
		}
		switch platformArg {
		case 1:
			fmt.Println("WELCOME TO THE PLATFORM OF CONTRACTS IDEA")
			// the platform of contracts idea is the idea of having an open platform with
			// various stakeholders and optimizing their game theoretic objectives. Specific
			// conttracts coud emulate the function of platforms themselves but in this case,
			// we handle all aspects of game theory within what we have and ecnourage people
			// to tak part in the system.
			// TODO: think of how to implement this in a nice way
		case 2:
			fmt.Println("WELCOME TO THE PLATFORM OF PLATFORMS IDEA")
			// the platform of platforms idea is similar to what we want in the opensolar
			// project where we can implement various partners as entities in the system
			// TODO: emulate different partners in the system that we have now.
		}
	*/
	// the memo field in stellar is particularly interesting for us wrt preserving state
	// as shown in the preliminary pdf example. We need to test out whether we can commit
	// more state in the memo field and find a way for third party clients to be able
	// to replicate this state so that this is verifiable and not just an entry of state
	// in the system.
	// TODO: think of a reasonable way to hash the current state of the system with
	// its hash in the memo field
	// TODO: write a receiver that can be run on the client (electricity generation thing)
	// which can relay the information out to us. We do need to create public and privatekey
	// pairs on the device, this is something that atonomi does well, so maybe talk to them
	// regarding this.
	// instead of fetching data after it passes through a 3rd party, it would be nice if we could
	// get the data and then pass it on to them since it has to interface with our smart contract
	// which interfaces with stellar. This is easier if we have a stellar client running on local,
	// but I think that would not be possible on a small device (or maybe too much work, idk)
	// does it need a remote stellar node running?

	platformPublicKey, platformSeed, err := StartPlatform()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("PLATFORM SEED IS: %s\n PLATFORM PUBLIC KEY IS: %s", platformSeed, platformPublicKey)
	// TODO: how much do we pay the investor? how does it work
	// Do we sell the REC created from the solar panels only to the investor? If so,
	// isn't that enough to propel investment in the solar contract itself?
	// TODO: need a server to run a public stellar node to test out stuff
	// change the API mapping
	// move current number of years metric to a separate package since that is
	// more suitable for a model like affordable housing.
	// look into what kind of data we get from the pi and checkout pi specific code
	// to see if we can get something from there.

	// TODO: so the idea would be to split the current PoC into two parts, one
	// focused on the platform of platform ideas and one on the platform of contracts
	// ideas. The current implementation that we have is focused more on the platform
	// of platforms idea, with ideas to integrate various partners at differnet stages
	// to use their input in some parts, but we could do it in an automated way with
	// assets and tokens for everything as well, which is similar to the platform
	// of contracts idea.
	// TODO: how do we emulate various partners? need to get input / have some stuff
	// that we presume they do and then fill in the rest.
	// also need to implement stages in contracts based on finalization
	// need to integrate ipfs into the workflow so that we can store copies of the
	// specific contract and then reference it when needed using the ipfs hash.
	// need to start working on the base contract that connects investors and
	// recipients and the investors and the platform. Need to transition automatically
	// and also cover breach scenarios in case the recipient doesn't pay for a specific
	// period of time
	// TODO: upgrade the RPC to fit in with recent changes
	// TODO: also write a Makefile so that its easy for people to get started with stuff
	// TODO: look into how flags are set and set flags on accounts - no documentation is around
	// regarding this for go, so idk
	fmt.Println("------------STELLAR HOUSE INVESTMENT CLI INTERFACE------------")

	// init stablecoin stuff
	err = stablecoin.InitStableCoin()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("------WHAT DO YOU WANT TO DO?------")
	fmt.Println("1. CREATE A NEW INVESTOR ACCOUNT")
	fmt.Println("2. CREATE A NEW RECIPIENT ACCOUNT")
	fmt.Println("3: ALREADY HAVE AN ACCOUNT")
	opt, err := scan.ScanForInt()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return
	}
	switch opt {
	case 1:
		err := NewInvestorPrompt()
		if err != nil {
			log.Println(err)
			return
		}
	case 2:
		err := NewRecipientPrompt()
		if err != nil {
			log.Println(err)
			return
		}
	default:
	}

	investor, recipient, contractor, isRecipient, isContractor, err := LoginPrompt()
	if err != nil {
		log.Println(err)
		return
	}
	// check if the user is a recipient here
	if isRecipient {
		// we already have the recipient, so no need to make a call to the database
		// have a for loop here with various options
		for {
			fmt.Println("------------RECIPIENT INTERFACE------------")
			fmt.Println("----CHOOSE ONE OF THE FOLLOWING OPTIONS----")
			fmt.Println("  1. Display all Open Projects (STAGE 3)")
			fmt.Println("  2. Display my Profile")
			fmt.Println("  3. Payback towards an Project (STAGE 6)")
			fmt.Println("  4. Exchange XLM for USD")
			fmt.Println("  5. Finalize a specific Project (STAGE 2->3)")
			fmt.Println("  6. Originate a specific Project (STAGE 0->1)")
			fmt.Println("  7. View all Projects (ALL STAGES)")
			fmt.Println("  8. View all Origin Projects (STAGE 1)")
			fmt.Println("  9. View All Balances ")
			fmt.Println("  10. Unlock Account")
			fmt.Println("  default: Exit")
			optI, err := scan.ScanForInt()
			if err != nil {
				fmt.Println("Couldn't read user input")
				break
			}
			switch optI {
			case 1:
				Stage3ProjectsDisplayPrompt()
				break
			case 2:
				PrintRecipient(recipient)
				break
			case 3:
				// TODO: migrate this to a contract model which is based off stages rather than using DBParams here
				log.Println(recipient.ReceivedSolarProjects)
				fmt.Println("WHICH PROJECT DO YOU WANT TO PAY BACK TOWARDS? (ENTER PROJECT NUMBER)")
				projectNumber, err := scan.ScanForInt()
				if err != nil {
					log.Println("INPUT NOT AN INTEGER, TRY AGAIN")
					continue
				}
				// check if we can get the project using the project number that we have here
				rtContract, err := solar.RetrieveProject(projectNumber)
				if err != nil {
					log.Println("Couldn't retrieve project, try again!")
					continue
				}
				// so we can retrieve the project using the project Index, nice
				PrintProject(rtContract)
				fmt.Println("HOW MUCH DO YOU WANT TO PAYBACK?")
				paybackAmount, err := scan.ScanForStringWithCheckI()
				if err != nil {
					log.Println(err)
					break
				}
				fmt.Printf(" DO YOU WANT TO CONFIRM THAT YOU WANT TO PAYBACK %s TOWARDS THIS PROJECT? (PRESS N IF YOU DON'T WANT TO)\n", paybackAmount)
				confirmOpt, err := scan.ScanForString()
				if err != nil {
					log.Println(err)
					break
				}
				if confirmOpt == "N" || confirmOpt == "n" {
					fmt.Println("YOU HAVE DECIDED TO CANCEL THE PAYBACK ORDER")
					break
				}
				fmt.Printf("PAYING BACK %s TOWARDS PROJECT NUMBER: %d\n", paybackAmount, rtContract.Params.Index) // use the rtContract.Params here instead of using projectNumber from long ago
				// now we need to call back the payback function to payback the asset
				// Here, we will simply payback the DEBTokens that was sent to us earlier
				if rtContract.Params.DebtAssetCode == "" {
					log.Fatal("Project not found")
				}

				err = solar.Payback(recipient.U.Index, rtContract.Params.Index, rtContract.Params.DebtAssetCode, platformPublicKey, platformSeed, paybackAmount, recipientSeed)
				// TODO: right now, the payback asset directly sends back, change
				if err != nil {
					log.Println("PAYBACK TX FAILED, PLEASE TRY AGAIN!")
					break
				}
				break
			case 4:
				log.Println("Enter the amount you want to convert into STABLEUSD")
				convAmount, err := scan.ScanForStringWithCheckF()
				if err != nil {
					log.Println(err)
					break
				}
				err = stablecoin.Exchange(recipient.U.PublicKey, recipientSeed, convAmount)
				if err != nil {
					log.Println(err)
					break
				}
				break
			case 5:
				allContracts, err := solar.RetrieveRecipientProjects(solar.ProposedProject, recipient.U.Index)
				if err != nil {
					log.Fatal(err)
				}
				PrintProjects(allContracts)

				fmt.Println("CHOOSE THE METRIC BY WHICH YOU WANT TO SELECT THE WINNING BID: ")
				fmt.Println("1. PRICE")
				fmt.Println("2. COMPLETION TIME (IN YEARS)")
				fmt.Println("3. SELECT MANUALLY")
				fmt.Println("ENTER YOUR CHOICE AS A NUMBER (1 / 2 / 3)")
				opt, err := scan.ScanForInt()
				if err != nil {
					fmt.Println("Couldn't read user input")
					break
				}
				// the length of the above two slices must be the same
				for _, elem := range allContracts {
					log.Println("======================================================================================")
					PrintProject(elem)
				}
				switch opt {
				case 1:
					fmt.Println("YOU'VE CHOSEN TO SELECT BY LEAST PRICE")
					// here we assume that the timeout period for the auction is up and that
					// price is the winning metric of a specific bid, like in traditional contract
					bestContract, err := solar.SelectContractByPrice(allContracts)
					if err != nil {
						log.Fatal(err)
					}
					log.Println("BEST CONTRACT IS: ")
					// we need the contractor who proposed this contract
					bestContract.SetFinalizedProject()
					PrintProject(bestContract)
					// now at this point, we need to mark this specific contract as completed.
					// do we set a flag? db entry? how do we do that
				case 2:
					fmt.Println("YOU'VE CHOSEN TO SELECT BY NUMBER OF YEARS")
					bestContract, err := solar.SelectContractByTime(allContracts)
					if err != nil {
						log.Fatal(err)
					}
					log.Println("BEST CONTRACT IS: ")
					bestContract.SetFinalizedProject()
					PrintProject(bestContract)
				case 3:
					for i, contract := range allContracts {
						log.Println("BEST BID CHOICE NUMBER: ", i)
						PrintProject(contract)
					}
					fmt.Println("ENTER YOUR OPTION AS A NUMBER")
					opt, err := scan.ScanForInt()
					if err != nil {
						log.Println(err)
						break
					}
					log.Println("BEST CONTRACT IS: ")
					// we need the contractor who proposed this contract
					allContracts[opt].SetFinalizedProject()
					PrintProject(allContracts[opt])
				default:
					break
				}
			case 6:
				fmt.Println("LIST OF ALL PRE ORIGIN PROJECTS BY ORIGINATORS (STAGE 0)")
				allMyProjects, err := solar.RetrieveProjectsAtStage(solar.PreOriginProject)
				if err != nil {
					log.Println(err)
					continue
				}
				PrintProjects(allMyProjects)
				fmt.Println("ENTER THE PROJECT INDEX")
				contractIndex, err := scan.ScanForInt()
				if err != nil {
					log.Println(err)
					continue
				}
				err = solar.PromoteStage0To1Project(contractIndex)
				if err != nil {
					log.Println(err)
					break
				}
			case 7:
				fmt.Println("PRINTING ALL PROJECTS: ")
				allContracts, err := solar.RetrieveAllProjects()
				if err != nil {
					log.Println(err)
					break
				}
				PrintProjects(allContracts)
			case 8:
				DisplayOriginProjects()
			case 9:
				BalanceDisplayPrompt(recipient.U.PublicKey)
			case 10:
				// need to unlock the recipient account
				seedpwd, err := scan.ScanRawPassword()
				if err != nil {
					log.Println(err)
					break
				}
				seed, err := wallet.DecryptSeed(recipient.U.EncryptedSeed, seedpwd)
				if err != nil {
					log.Println(err)
					break
				}
				recipientSeed = seed
				log.Println(" Seed successfully unlocked")
			default: // this default is for the larger switch case
				ExitPrompt()
			}
		}
		return
	} else if isContractor {
		var optI int
		if contractor.Contractor {
			optI = 1
		} else if contractor.Originator {
			optI = 2
		}
		switch optI {
		case 1:
			fmt.Println("-------------WELCOME BACK CONTRACTOR-------------")
			for {
				fmt.Println("WHAT WOULD YOU LIKE TO DO?")
				fmt.Println("  1. VIEW ALL ORIGINATED (STAGE 1) PROJECTS")
				fmt.Println("  2. VIEW PROFILE")
				fmt.Println("  3. CREATE A PROPOSED (STAGE 2) PROJECT")
				fmt.Println("  4. VIEW ALL MY PROPOSED (STAGE 2) PROJECTS")
				optI, err := scan.ScanForInt()
				if err != nil {
					log.Println(err)
					continue
				}
				switch optI {
				case 1:
					DisplayOriginProjects()
				case 2:
					PrintEntity(contractor)
				case 3:
					fmt.Println("YOU HAVE CHOSEN TO CREATE A NEW PROPOSED PROJECT")
					err = ProposeContractPrompt(&contractor)
					if err != nil {
						log.Println(err)
						continue
					}
				case 4:
					fmt.Println("LIST OF ALL PROPOSED CONTRACTS BY ME: ")
					allMyProjects, err := solar.RetrieveContractorProjects(solar.ProposedProject, contractor.U.Index)
					if err != nil {
						log.Println(err)
						continue
					}
					PrintProjects(allMyProjects)
				}
			}
		case 2:
			fmt.Println("-------------WELCOME BACK ORIGINATOR-------------")
			for {
				fmt.Println("WHAT WOULD YOU LIKE TO DO?")
				fmt.Println("  1. PROPOSE A PRE-ORIGIN (STAGE 0) PROJECT TO A RECIPIENT")
				fmt.Println("  2. VIEW PROFILE")
				fmt.Println("  3. VIEW ALL MY PRE-ORIGINATED (STAGE 0) PROJECTS")
				fmt.Println("  4. VIEW ALL MY ORIGINATED (STAGE 1) PROJECTS")
				optI, err = scan.ScanForInt()
				if err != nil {
					log.Println(err)
					continue
				}
				switch optI {
				case 1:
					err := OriginContractPrompt(&contractor)
					if err != nil {
						fmt.Println("RETURNING BACK TO THE MAIN LOOP: ", err)
						continue
					}
				case 2:
					PrintEntity(contractor)
				case 3:
					allMyProjects, err := solar.RetrieveOriginatorProjects(solar.PreOriginProject, contractor.U.Index)
					if err != nil {
						fmt.Println("RETURNING BACK TO THE MAIN LOOP: ", err)
						continue
					}
					PrintProjects(allMyProjects)
				case 4:
					allMyProjects, err := solar.RetrieveOriginatorProjects(solar.OriginProject, contractor.U.Index)
					if err != nil {
						fmt.Println("RETURNING BACK TO THE MAIN LOOP: ", err)
						continue
					}
					PrintProjects(allMyProjects)
				default:
					ExitPrompt()
				}
			}
		}
		PrintEntity(contractor)
	} else {
		// User is an investor
		for {
			// Main investor loop
			fmt.Println("------------INVESTOR INTERFACE------------")
			fmt.Println("----CHOOSE ONE OF THE FOLLOWING OPTIONS----")
			fmt.Println("  1. Display all Open Projects (STAGE 3)")
			fmt.Println("  2. Display my Profile")
			fmt.Println("  3. Invest in an Project (STAGE 3)")
			fmt.Println("  4. Display All Balances")
			fmt.Println("  5. Exchange XLM for USD")
			fmt.Println("  6. Display all Origin (STAGE 1) Projects")
			fmt.Println("  7. Vote towards a specific proposed project (STAGE 2)")
			fmt.Println("  8. Get ipfs hash of a contract")
			fmt.Println("  9. Display all Funded Projects")
			fmt.Println("  10. Unlock account")
			fmt.Println("  default: Exit")
			optI, err := scan.ScanForInt()
			if err != nil {
				fmt.Println("Couldn't read user input")
				break
			}
			switch optI {
			case 1:
				Stage3ProjectsDisplayPrompt()
				break
			case 2:
				PrintInvestor(investor)
				break
			case 3:
				fmt.Println("----WHICH PROJECT DO YOU WANT TO INVEST IN? (ENTER ORDER NUMBER WITHOUT SPACES)----")
				oNumber, err := scan.ScanForInt()
				if err != nil {
					fmt.Println("Couldn't read user input")
					break
				}
				// now the user has decided to invest in the asset with index uInput
				// we need to retrieve the project and ask for confirmation
				solarProject, err := solar.RetrieveProject(oNumber)
				if err != nil {
					log.Println("Couldn't retrieve project, try again!")
					continue
				}
				PrintProject(solarProject)
				fmt.Println(" HOW MUCH DO YOU WANT TO INVEST?")
				investmentAmount, err := scan.ScanForStringWithCheckI()
				if err != nil {
					log.Println(err)
					break
				}
				fmt.Println(" DO YOU WANT TO CONFIRM THIS ORDER? (PRESS N IF YOU DON'T WANT TO)")
				confirmOpt, err := scan.ScanForString()
				if err != nil {
					log.Println(err)
					break
				}
				if confirmOpt == "N" || confirmOpt == "n" {
					fmt.Println("YOU HAVE DECIDED TO CANCEL THIS ORDER")
					break
				}

				err = platform.RefillPlatform(platformPublicKey)
				if err != nil {
					log.Println(err)
					break
				}
				fmt.Printf("Platform seed is: %s and platform's publicKey is %s", platformSeed, platformPublicKey)
				err = xlm.RefillAccount(investor.U.PublicKey, platformSeed)
				if err != nil {
					log.Println(err)
					break
				}
				recipient := solarProject.ProjectRecipient
				// from here on, reference recipient
				err = xlm.RefillAccount(recipient.U.PublicKey, platformSeed)
				if err != nil {
					log.Println(err)
					break
				}

				platformBalance, err := xlm.GetNativeBalance(platformPublicKey)
				if err != nil {
					log.Fatal(err)
				}

				// need the recipient's seed here as well
				// need to unlock the recipient account
				fmt.Println("ENTER THE RECIPIENT'S SEED PASSWORD")
				// ideally we should ask the recipient for confirmation in case he wants to re4ceived the money or something
				seedpwd, err := scan.ScanRawPassword()
				if err != nil {
					log.Println(err)
					break
				}
				seed, err := wallet.DecryptSeed(recipient.U.EncryptedSeed, seedpwd)
				if err != nil {
					log.Println(err)
					break
				}
				recipientSeed = seed
				log.Println(" Seed successfully unlocked")
				log.Println("Platform's updated balance is: ", platformBalance)
				log.Println("The investor's public key and private key are: ", investor.U.PublicKey, " ", investorSeed)
				log.Println("The recipient's public key and private key are: ", recipient.U.PublicKey, " ", recipientSeed)
				// so now we have three entities setup, so we create the assets and invest in them
				cProject, err := solar.InvestInProject(solarProject.Params.Index, platformPublicKey, platformSeed, investor.U.Index, recipient.U.Index, investmentAmount, investorSeed, recipientSeed) // assume payback period is 5
				if err != nil {
					log.Println(err)
				} else {
					fmt.Println("YOUR PROJECT INVESTMENT HAS BEEN CONFIRMED: ")
					PrintProject(cProject)
					fmt.Println("PLEASE CHECK A BLOCKCHAIN EXPLORER TO CONFIRM BALANCES: ")
					fmt.Println("https://testnet.steexp.com/account/" + investor.U.PublicKey + "#balances")
				}
				break
			case 4:
				BalanceDisplayPrompt(investor.U.PublicKey)
				break
			case 5:
				// this should be expanded in the future to make use of the inbuilt DEX
				// on stellar (checkout stellarterm)
				log.Println("Enter the amount you want to convert into STABLEUSD")
				// this would also mean that you need to check whether we have the balance
				// here and then proceed further
				convAmount, err := scan.ScanForStringWithCheckF()
				if err != nil {
					log.Println(err)
					return
				}
				// maybe don't trust asset again when you've trusted it already? check if that's
				// possible and save on the tx fee for a single transaction. But I guess its
				// difficult to retrieve trustlines, so we'll go ahead with it
				hash, err := assets.TrustAsset(stablecoin.StableUSD, consts.StablecoinTrustLimit, investor.U.PublicKey, investorSeed)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("tx hash for trusting stableUSD: ", hash)
				// now send coins across and see if our tracker detects it
				_, hash, err = xlm.SendXLM(stablecoin.PublicKey, convAmount, investorSeed, "sending xlm to bootstrap")
				if err != nil {
					log.Fatal(err)
				}

				log.Println("tx hash for sent xlm: ", hash, "pubkey: ", investor.U.PublicKey)
				break
			case 6:
				DisplayOriginProjects()
			case 7:
				// this is the case where an investor can vote on a particular proposed project
				// ie stage 2 projects for the recipient to have an understanding about
				// which contracts are popular and can receive more investor money
				fmt.Println("LIST OF ALL PROPOSED ORDERS: ")
				allProposedProjects, err := solar.RetrieveProjectsAtStage(solar.ProposedProject)
				if err != nil {
					log.Println(err)
					break
				}
				PrintProjects(allProposedProjects)
				fmt.Println("WHICH CONTRACT DO YOU WANT TO VOTE TOWARDS?")
				vote, err := scan.ScanForInt()
				if err != nil {
					log.Println(err)
					break
				}
				fmt.Println("YOUR AVAILABLE VOTING BALANCE IS: ", investor.VotingBalance)
				fmt.Println("HOW MANY VOTES DO YOU WANT TO DELEGATE TOWARDS THIS ORDER?")
				votes, err := scan.ScanForInt()
				if err != nil {
					log.Println(err)
					break
				}
				err = solar.VoteTowardsProposedProject(investor.U.Index, votes, vote)
				if err != nil {
					log.Println(err)
					break
				}
				log.Println("You have voted for contract number: ", vote)
			case 8:
				fmt.Println("WELCOME TO THE IPFS HASHING INTERFACE")
				fmt.Println("ENTER THE STRING THAT YOU WOULD LIKE THE IPFS HASH FOR")
				// the UI should ideally have a menu that asks the user for a file and then
				// produces the hash of it. In this case, we shall use a sample ipfs file
				// and then hash it.
				// this uses the platform's ipfs key though, not the user's. If the user
				// wants to serve his own ipfs files, he is better off running a client on
				// his own
				hashString, err := scan.ScanForString()
				if err != nil {
					fmt.Println("Couldn't read user input, going back to the main loop")
				}
				hash, err := ipfs.AddStringToIpfs(hashString)
				if err != nil {
					fmt.Println("Couldn't hash user input, exiting to main menu", err)
					break
				}
				hashCheck, err := ipfs.GetStringFromIpfs(hash)
				if err != nil || hashCheck != hashString {
					fmt.Println("Hashed strings and retrieved strings don't match, don't use this hash!", err)
					break
				}
				// don't print this hash unless we can decrypt it and be sure that it behaves as expected
				log.Println("THE HASH OF THE PROVIDED STRING IS: ", hash)
				// try to retrieve the string back from ipfs and check if it works correctly
			case 9:
				fmt.Println("LIST OF ALL FUNDED PROJECTS: ")
				allFundedProjects, err := solar.RetrieveProjectsAtStage(solar.FundedProject)
				if err != nil {
					fmt.Println(err)
					break
				}
				PrintProjects(allFundedProjects)
			case 10:
				// need to unlock the recipient account
				seedpwd, err := scan.ScanRawPassword()
				if err != nil {
					log.Println(err)
					break
				}
				seed, err := wallet.DecryptSeed(investor.U.EncryptedSeed, seedpwd)
				if err != nil {
					log.Println(err)
					break
				}
				investorSeed = seed
				log.Println(" Seed successfully unlocked: ", seed)
			default:
				ExitPrompt()
			} // end of switch
		}
		// it should never arrive here
		return
	}
}

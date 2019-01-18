package solar

import (
	"fmt"
	"log"

	assets "github.com/OpenFinancing/openfinancing/assets"
	consts "github.com/OpenFinancing/openfinancing/consts"
	database "github.com/OpenFinancing/openfinancing/database"
	utils "github.com/OpenFinancing/openfinancing/utils"
	"github.com/stellar/go/build"
)

// this file does not contain any tests associated with it right now. In the future,
// once we have a robust frontend, we can modify the CLI interface to act as a test
// for this file

// InvestInProject invests in a particular solar project given required parameters
func InvestInProject(projIndex int, issuerPubkey string, issuerSeed string, invIndex int, recpIndex int, invAmountS string, invSeed string, recpSeed string) (Project, error) {
	var err error

	project, err := RetrieveProject(projIndex)
	if err != nil {
		return project, err
	}
	// invest only in integer values as of now, TODO: change to float
	invAmount := utils.StoI(invAmountS)
	// check if investment amount is greater than or equal to the project requirements
	if invAmount > project.Params.TotalValue-project.Params.MoneyRaised {
		fmt.Println("User is trying to invest more than what is needed, print and exit")
		return project, fmt.Errorf("User is trying to invest more than what is needed, print and exit")
	}
	investor, err := database.RetrieveInvestor(invIndex)
	if err != nil {
		return project, err
	}
	// we should check here whether the investor has enough USDTokens in project to be
	// able to ivnest in the asset
	if !investor.CanInvest(investor.U.PublicKey, invAmountS) {
		log.Println("Investor has less balance than what is required to ivnest in this asset")
		return project, err
	}
	recipient, err := database.RetrieveRecipient(recpIndex)
	if err != nil {
		return project, err
	}

	// user has decided to invest in a part of the project (don't know if full yet)
	// so if there has been no token codes assigned yet, we need to create them and
	// assign them here
	// you can retrieve these anywhere since the metadata will most likely be unique
	assetName := assets.AssetID(project.Params.Metadata)
	if project.Params.InvestorAssetCode == "" {
		// this person is the first investor, set the investor token name
		project.Params.InvestorAssetCode = assets.AssetID(consts.InvestorAssetPrefix + assetName) // set the investeor code
		_ = assets.CreateAsset(project.Params.InvestorAssetCode, issuerPubkey)                    // create the asset itself, since it would not have bene created earlier
	}
	var InvAsset build.Asset
	InvAsset.Code = project.Params.InvestorAssetCode
	InvAsset.Issuer = issuerPubkey
	// InvAsset is not a native token, so don't set native to true
	// make investor trust the asset, trustlimit is upto the value of the project
	txHash, err := assets.TrustAsset(InvAsset, utils.ItoS(project.Params.TotalValue), investor.U.PublicKey, invSeed)
	if err != nil {
		return project, err
	}
	log.Println("Investor trusted asset: ", InvAsset.Code, " tx hash: ", txHash)
	_, txHash, err = assets.SendAssetFromIssuer(InvAsset.Code, investor.U.PublicKey, invAmountS, issuerSeed, issuerPubkey)
	if err != nil {
		return project, err
	}
	log.Printf("Sent InvAsset %s to investor %s with txhash %s", InvAsset.Code, investor.U.PublicKey, txHash)
	// investor asset sent, update project.Params's BalLeft
	project.Params.MoneyRaised += invAmount
	fmt.Println("Updating investor to handle invested amounts and assets")
	investor.AmountInvested += float64(invAmount)
	investor.InvestedSolarProjects = append(investor.InvestedSolarProjects, project.Params.DebtAssetCode)
	err = investor.Save() // save investor creds now that we're done
	if err != nil {
		return project, err
	}
	fmt.Println("Updated investor database")
	// append the investor class to the list of project investors
	// if the same investor has invested twice, he will appear twice
	// can be resolved on the UI side by requiring unique, so not doing that here
	project.ProjectInvestors = append(project.ProjectInvestors, investor)
	if project.Params.MoneyRaised == project.Params.TotalValue {
		// this project covers up the amount nedeed for the project, so set the DebtAssetCode
		// and PaybackAssetCodes, generate them and give to the recipient
		project.Params.DebtAssetCode = assets.AssetID(consts.DebtAssetPrefix + assetName)
		project.Params.PaybackAssetCode = assets.AssetID(consts.PaybackAssetPrefix + assetName)
		DebtAsset := assets.CreateAsset(project.Params.DebtAssetCode, issuerPubkey)
		PaybackAsset := assets.CreateAsset(project.Params.PaybackAssetCode, issuerPubkey)
		// and the school needs to trust me only for paybackTokens amount of PB tokens
		pbAmtTrust := utils.ItoS(project.Params.Years * 12 * 2) // two way exchange possible, to account for errors
		txHash, err = assets.TrustAsset(PaybackAsset, pbAmtTrust, recipient.U.PublicKey, recpSeed)
		if err != nil {
			return project, err
		}
		log.Println("Recipient Trusted Payback asset: ", PaybackAsset.Code, " tx hash: ", txHash)
		txHash, err = assets.TrustAsset(DebtAsset, utils.ItoS(project.Params.TotalValue*2), recipient.U.PublicKey, recpSeed) // since debt = invested amount
		// *2 is for sending the amount back
		if err != nil {
			return project, err
		}
		log.Println("Recipient Trusted Debt asset: ", DebtAsset.Code, " tx hash: ", txHash)
		_, txHash, err = assets.SendAssetFromIssuer(project.Params.DebtAssetCode, recipient.U.PublicKey, utils.ItoS(project.Params.TotalValue), issuerSeed, issuerPubkey) // same amount as debt
		if err != nil {
			return project, err
		}
		log.Printf("Sent DebtAsset to recipient %s with txhash %s", recipient.U.PublicKey, txHash)
		project.Params.BalLeft = float64(project.Params.TotalValue)
		recipient.ReceivedSolarProjects = append(recipient.ReceivedSolarProjects, project.Params.DebtAssetCode)
		project.ProjectRecipient = recipient // need to udpate project.Params each time recipient is mutated
		// only here does the recipient part change, so update it only here
		// TODO: keep note of who all invested in this asset (even though it should be
		// easy to get that from the blockchain)
		if project.Params.DebtAssetCode == "" {
			return project, fmt.Errorf("Empty debt asset code")
		}
		err = recipient.Save()
		if err != nil {
			return project, err
		}
		project.Stage = FundedProject // set funded project stage
		err = project.Save()
		if err != nil {
			log.Println("Couldn't insert project")
			return project, err
		}
		fmt.Println("Updated recipient bucket")
		return project, nil
	}
	// update the project finally now that we have updated other databases
	err = project.Save()
	return project, err
}

// sendPaybackAsset sends back the solar pb asset from the issuer to the recipient
func sendPaybackAsset(project Project, destination string, amount string, platformSeed string, platformPubkey string) error {
	// need to calculate how much PaybackAsset we need to send back.
	amountS := project.CalculatePayback(amount)
	_, txHash, err := assets.SendAssetFromIssuer(project.Params.PaybackAssetCode, destination, amountS, platformSeed, platformPubkey)
	log.Println("TXHASH for payback is: ", txHash)
	return err
}

package solar

import (
	"fmt"
	"log"

	database "github.com/OpenFinancing/openfinancing/database"
	oracle "github.com/OpenFinancing/openfinancing/oracle"
	utils "github.com/OpenFinancing/openfinancing/utils"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
)

// Payback is called when the receiver of the DEBToken wants to pay a fixed amount
// of money back to the issuer of the DEBTokens. One way to imagine this would be
// like an electricity bill, something that people pay monthly but only that in this
// case, the electricity is free, so they pay directly towards the solar panels.
// The process of Payback roughly involves the followign steps:
// 1. Pay the issuer in DEBTokens with whatever amount desired.
// The oracle price of
// electricity cost is a lower bound (since the government would not like it if people
// default on their payments). (MW: Explain this lower bound and default issue more)
// Anything below the lower bound gets a warning in
// project for people to pay more, we could also have a threshold mechanism that says
// if a person constantly defaults for more than half the owed amount for three
// consecutive months, we sell power directly to the grid. THis could also be used
// for a rating system, where the frontend UI can have a rating based on whether
// the recipient has defaulted or not in the past.
// 2. The receiver checks whether the amount is greater than Oracle Threshold and
// if so, sends back PBTokens, which stand for the month equivalent of payments.
// eg. the school has opted for a 5 year payback period, the school owes the issuer
// 60 PBTokens and the issuer sends back 1PBToken every month if the school pays
// invested_amount/60 DEBTokens back to the issuer
// 3. The recipient checks whether the PBTokens received correlate to the amount
// that it sent and if not, raises the dispute since the forward DEBToken payment
// is on chain and resolves the dispute itself using existing off chain legal frameworks
// (issued bonds, agreements, etc)
func Payback(recpIndex int, projIndex int, assetName string, issuerPubkey string, issuerSeed string, amount string, seed string) error {
	recipient, err := database.RetrieveRecipient(recpIndex)
	if err != nil {
		return err
	}
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return err
	}
	// once we have the stablecoin here, we can remove the assetName
	StableBalance, err := xlm.GetAssetBalance(recipient.U.PublicKey, "STABLEUSD")
	// checks for the stablecoin asset
	if err != nil {
		log.Println("YOU HAVE NO STABLECOIN BALANCE, PLEASE REFILL ACCOUNT")
		return err
	}

	DEBAssetBalance, err := xlm.GetAssetBalance(recipient.U.PublicKey, assetName)
	if err != nil {
		log.Println("Don't have the debt asset in possession")
		return err
	}

	if utils.StoF(amount) > utils.StoF(StableBalance) {
		// check whether the recipient has enough StableUSD tokens in project to make
		// this happen
		log.Println("YOU CAN'T SEND AN AMOUNT MORE THAN WHAT YOU HAVE")
		return fmt.Errorf("YOU CAN'T SEND AN AMOUNT MORE THAN WHAT YOU HAVE")
	}
	// check balance in DEBAssetCode anmd
	monthlyBill := oracle.MonthlyBill()
	if err != nil {
		log.Println("Unable to fetch oracle price, exiting")
		return err
	}

	log.Println("Retrieved average price from oracle: ", monthlyBill)
	// the oracle needs to know the assetName so that it can find the other details
	// about this asset from the db. This should run on the server side and must
	// be split when we do run client side stuff.
	// hardcode for now, need to add the oracle here so that we
	// can do this dynamically
	// send amount worth DEBTokens back to issuer
	confHeight, txHash, err := recipient.SendAssetToIssuer(assetName, issuerPubkey, amount, seed)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Paid debt amount: ", amount, " back to issuer, tx hash: ", txHash, " ", confHeight)
	// log.Println("Checking balance to see if our account was debited")
	newBalance, err := xlm.GetAssetBalance(recipient.U.PublicKey, assetName)
	if err != nil {
		return err
	}

	newBalanceFloat := utils.StoF(newBalance)
	DEBAssetBalanceFloat := utils.StoF(DEBAssetBalance)
	mBillFloat := utils.StoF(monthlyBill)

	paidAmount := DEBAssetBalanceFloat - newBalanceFloat
	log.Println("Old Balance: ", DEBAssetBalanceFloat, "New Balance: ", newBalanceFloat, "Paid: ", paidAmount, "Bill Amount: ", mBillFloat)

	// would be nice to take some additional action like sending a notification or
	// something to investors or to the email address given so that everyone is made
	// aware of this and there's data transparency

	if paidAmount < mBillFloat {
		log.Println("Amount paid is less than amount required, balance not updating, please make sure to cover this next time")
	} else if paidAmount > mBillFloat {
		log.Println("You've chosen to pay more than what is required for this month. Adjusting payback period accordingly")
	} else {
		log.Println("You've paid exactly what is required for this month. Payback period remains as usual")
	}
	// we need to update the database here
	// no need to retrieve this project again because we have it already
	project.Params.BalLeft -= paidAmount
	project.Params.DateLastPaid = utils.Timestamp()
	if project.Params.BalLeft == 0 {
		log.Println("YOU HAVE PAID OFF THIS ASSET, TRANSFERRING OWNERSHIP OF ASSET TO YOU")
		// don't delete the asset from the received assets list, we still need it so
		// that we can look back and find out hwo many assets this particular
		// enttiy has been invested in, have a leaderboard kind of thing, etc.
		project.Stage = 7
		// we should call neighbourly or some ohter partner here to transfer assets
		// using the bond they provide us with
		// the nice part here is that the recipient can not pay off more than what is
		// invested because the trustline will not allow such an incident to happen
	}
	// balLeft must be updated on the server side and can be challenged easily
	// if there's some discrepancy since the tx's are on the blockchain
	err = sendPaybackAsset(project, recipient.U.PublicKey, amount, issuerSeed, issuerPubkey)
	if err != nil {
		log.Println("Sending back payback asset failed, please try again!", err)
		return err
	}

	err = project.updateRecipient(recipient)
	if err != nil {
		return err
	}

	err = project.Save()
	if err != nil {
		return err
	}
	return err
}

// CalculatePayback is a TODO function that should simply sum the PBToken
// balance and then return them to the frontend UI for a nice display
func (project Project) CalculatePayback(amount string) string {
	// the idea is that we should be able to pass an assetId to this function
	// and it must calculate how much time we have left for payback. For this example
	// until we do the db stuff, lets pass a few params (although this could be done
	// separately as well).
	// TODO: this functon needs to be the payback function
	// Consider not only the PBTokens amount, but also the rate of solar generation or consumption per year
	// (Eg. PBTokens are 10'000 and the rate/yrs is 2000, thus 5 years is the estimated bond maturity and payback remaining time)
	amountF := utils.StoF(amount)
	amountPB := (amountF / float64(project.Params.TotalValue)) * float64(project.Params.Years*12)
	amountPBString := utils.FtoS(amountPB)
	return amountPBString
}

package bonds

import (
	"encoding/json"
	"fmt"
	"log"

	assets "github.com/YaleOpenLab/openx/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	utils "github.com/YaleOpenLab/openx/utils"
	"github.com/boltdb/bolt"
)

// the coop struct uses the same base params as the bond model
type Coop struct {
	Params         BondCoopParams
	UnitsSold      int
	TotalAmount    float64
	TypeOfUnit     string
	MonthlyPayment float64
	Residents      []database.Investor
}

// NewCoop returns a new living coop and automatically saves it
func NewCoop(mdate string, mrights string, stype string, intrate float64, rating string,
	bIssuer string, uWriter string, totalAmount float64, typeOfUnit string, monthlyPayment float64,
	title string, location string, description string) (Coop, error) {
	var cCoop Coop
	cCoop.Params = newParams(mdate, mrights, stype, intrate, rating, bIssuer, uWriter, title, location, description)
	x, err := RetrieveAllCoops()
	if err != nil {
		return cCoop, err
	}

	cCoop.Params.Index = len(x) + 1
	cCoop.UnitsSold = 0
	cCoop.TotalAmount = totalAmount
	cCoop.TypeOfUnit = typeOfUnit
	cCoop.MonthlyPayment = monthlyPayment
	err = cCoop.Save()
	return cCoop, err
}

func (a *Coop) Save() error {
	db, err := database.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.CoopBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.ItoB(a.Params.Index)), encoded)
	})
	return err
}

// RetrieveAllCoops gets a list of all User in the database
func RetrieveAllCoops() ([]Coop, error) {
	var arr []Coop
	db, err := database.OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.CoopBucket)
		for i := 1; ; i++ {
			var rCoop Coop
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rCoop)
			if err != nil {
				return err
			}
			arr = append(arr, rCoop)
		}
		return nil
	})
	return arr, err
}

// RetrieveCoop retrieves a specifi coop from the database
func RetrieveCoop(key int) (Coop, error) {
	var bond Coop
	db, err := database.OpenDB()
	if err != nil {
		return bond, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.CoopBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return fmt.Errorf("Retrieved Coop nil")
		}
		return json.Unmarshal(x, &bond)
	})
	return bond, err
}

// for the demo, the publickey and seed must be hardcoded and  given as a binary I guess
// or worse, hardcode the seed and pubkey in the functions themselves
// TODO: add the recipient's role here, whether to give him an asset or do nothing
func (a *Coop) Invest(issuerPublicKey string, issuerSeed string, investor *database.Investor,
	investmentAmountS string, investorSeed string) error {
	// we want to invest in this specific bond
	var err error
	investmentAmount := utils.StoI(investmentAmountS)
	// check if investment amount is greater than the cost of a unit
	if float64(investmentAmount) > a.MonthlyPayment || float64(investmentAmount) < a.MonthlyPayment {
		fmt.Println("You are trying to invest more or less than a month's payment")
		return fmt.Errorf("You are trying to invest more or less than a month's payment")
	}
	assetName := assets.AssetID(a.Params.MaturationDate + a.Params.SecurityType + a.Params.Rating + a.Params.BondIssuer) // get a unique assetID

	if a.Params.InvestorAssetCode == "" {
		// this person is the first investor, set the investor token name
		InvestorAssetCode := assets.AssetID(consts.CoopAssetPrefix + assetName)
		a.Params.InvestorAssetCode = InvestorAssetCode             // set the investeor code
		_ = assets.CreateAsset(InvestorAssetCode, issuerPublicKey) // create the asset itself, since it would not have bene created earlier
	}
	/*
		dont check stableUSD balance for now
		if !investor.CanInvest(investor.U.PublicKey, investmentAmountS) {
			log.Println("Investor has less balance than what is required to ivnest in this asset")
			return a, err
		}
	*/
	// make investor trust the asset that we provide
	txHash, err := assets.TrustAsset(a.Params.InvestorAssetCode, issuerPublicKey, utils.FtoS(a.TotalAmount), investor.U.PublicKey, investorSeed)
	// trust upto the total value of the asset
	if err != nil {
		return err
	}
	log.Println("Investor trusted asset: ", a.Params.InvestorAssetCode, " tx hash: ", txHash)
	log.Println("Sending INVAsset: ", a.Params.InvestorAssetCode, "for: ", investmentAmount)
	_, txHash, err = assets.SendAssetFromIssuer(a.Params.InvestorAssetCode, investor.U.PublicKey, investmentAmountS, issuerSeed, issuerPublicKey)
	if err != nil {
		return err
	}
	log.Printf("Sent INVAsset %s to investor %s with txhash %s", a.Params.InvestorAssetCode, investor.U.PublicKey, txHash)
	// investor asset sent, update a.Params's BalLeft
	a.UnitsSold += 1
	investor.AmountInvested += float64(investmentAmount)
	investor.InvestedCoops = append(investor.InvestedCoops, a.Params.InvestorAssetCode)
	err = investor.Save() // save investor creds now that we're done
	if err != nil {
		return err
	}
	a.Residents = append(a.Residents, *investor)
	err = a.Save()
	return err
}

package bonds

import (
	"encoding/json"
	"fmt"
	"log"

	assets "github.com/OpenFinancing/openfinancing/assets"
	consts "github.com/OpenFinancing/openfinancing/consts"
	database "github.com/OpenFinancing/openfinancing/database"
	utils "github.com/OpenFinancing/openfinancing/utils"
	"github.com/boltdb/bolt"
	"github.com/stellar/go/build"
)

type ConstructionBond struct {
	Params BondCoopParams
	// common set of params that we need for openfinancing
	AmountRaised   float64
	CostOfUnit     float64
	InstrumentType string
	NoOfUnits      int
	Tax            string
	Investors      []database.Investor
	RecipientIndex int
}

func newParams(mdate string, mrights string, stype string, intrate float64, rating string,
	bIssuer string, uWriter string, title string, location string, description string) BondCoopParams {
	var rParams BondCoopParams
	rParams.MaturationDate = mdate
	rParams.MemberRights = mrights
	rParams.SecurityType = stype
	rParams.InterestRate = intrate
	rParams.Rating = rating
	rParams.BondIssuer = bIssuer
	rParams.Underwriter = uWriter
	rParams.Title = title
	rParams.Location = location
	rParams.Description = description
	rParams.DateInitiated = utils.Timestamp()
	return rParams
}

func NewBond(mdate string, mrights string, stype string, intrate float64, rating string,
	bIssuer string, uWriter string, unitCost float64, itype string, nUnits int, tax string, recIndex int,
	title string, location string, description string) (ConstructionBond, error) {
	var cBond ConstructionBond
	cBond.Params = newParams(mdate, mrights, stype, intrate, rating, bIssuer, uWriter, title, location, description)
	x, err := RetrieveAllBonds()
	if err != nil {
		return cBond, err
	}

	cBond.Params.Index = len(x) + 1
	cBond.CostOfUnit = unitCost
	cBond.InstrumentType = itype
	cBond.NoOfUnits = nUnits
	cBond.Tax = tax
	cBond.RecipientIndex = recIndex
	err = cBond.Save()
	return cBond, err
}

func (a *ConstructionBond) Save() error {
	db, err := database.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.BondBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.ItoB(a.Params.Index)), encoded)
	})
	return err
}

// RetrieveAllBonds gets a list of all User in the database
func RetrieveAllBonds() ([]ConstructionBond, error) {
	var arr []ConstructionBond
	db, err := database.OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.BondBucket)
		for i := 1; ; i++ {
			var rBond ConstructionBond
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rBond)
			if err != nil {
				return err
			}
			arr = append(arr, rBond)
		}
		return nil
	})
	return arr, err
}

func RetrieveBond(key int) (ConstructionBond, error) {
	var bond ConstructionBond
	db, err := database.OpenDB()
	if err != nil {
		return bond, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.BondBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return nil
		}
		return json.Unmarshal(x, &bond)
	})
	return bond, err
}

// for the demo, the publickey and seed must be hardcoded and  given as a binary I guess
// or worse, hardcode the seed and pubkey in the functions themselves
func (a *ConstructionBond) Invest(issuerPublicKey string, issuerSeed string, investor *database.Investor,
	recipient *database.Recipient, investmentAmountS string, investorSeed string, recipientSeed string) error {
	// we want to invest in this specific bond
	var err error
	investmentAmount := utils.StoI(investmentAmountS)
	// check if investment amount is greater than the cost of a unit
	if float64(investmentAmount) > a.CostOfUnit {
		fmt.Println("You are trying to invest more than a unit's cost, do you want to invest in two units?")
		return fmt.Errorf("You are trying to invest more than a unit's cost, do you want to invest in two units?")
	}
	assetName := assets.AssetID(a.Params.MaturationDate + a.Params.SecurityType + a.Params.Rating + a.Params.BondIssuer) // get a unique assetID

	if a.Params.INVAssetCode == "" {
		// this person is the first investor, set the investor token name
		INVAssetCode := assets.AssetID(consts.BondAssetPrefix + assetName)
		a.Params.INVAssetCode = INVAssetCode                  // set the investeor code
		_ = assets.CreateAsset(INVAssetCode, issuerPublicKey) // create the asset itself, since it would not have bene created earlier
	}
	/*
		dont check stableUSD balance for now
		if !investor.CanInvest(investor.U.PublicKey, investmentAmountS) {
			log.Println("Investor has less balance than what is required to ivnest in this asset")
			return a, err
		}
	*/
	var INVAsset build.Asset
	INVAsset.Code = a.Params.INVAssetCode
	INVAsset.Issuer = issuerPublicKey
	// make in v estor trust the asset that we provide
	txHash, err := assets.TrustAsset(INVAsset, utils.FtoS(a.CostOfUnit*float64(a.NoOfUnits)), investor.U.PublicKey, investorSeed)
	// trust upto the total value of the asset
	if err != nil {
		return err
	}
	log.Println("Investor trusted asset: ", INVAsset.Code, " tx hash: ", txHash)
	log.Println("Sending INVAsset: ", INVAsset.Code, "for: ", investmentAmount)
	_, txHash, err = assets.SendAssetFromIssuer(INVAsset.Code, investor.U.PublicKey, investmentAmountS, issuerSeed, issuerPublicKey)
	if err != nil {
		return err
	}
	log.Printf("Sent INVAsset %s to investor %s with txhash %s", INVAsset.Code, investor.U.PublicKey, txHash)
	// investor asset sent, update a.Params's BalLeft
	a.AmountRaised += float64(investmentAmount)
	investor.AmountInvested += float64(investmentAmount)
	investor.InvestedBonds = append(investor.InvestedBonds, a.Params.INVAssetCode)
	err = investor.Save() // save investor creds now that we're done
	if err != nil {
		return err
	}
	a.Investors = append(a.Investors, *investor)
	err = a.Save()
	return err
}

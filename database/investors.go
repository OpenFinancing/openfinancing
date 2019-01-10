package database

// contains the WIP Investor struct which will be stored in a separate bucket
import (
	"encoding/json"
	"fmt"
	"log"

	scan "github.com/OpenFinancing/openfinancing/scan"
	utils "github.com/OpenFinancing/openfinancing/utils"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
	"github.com/boltdb/bolt"
	"github.com/stellar/go/build"
)

// the investor struct contains all the investor details such as
// public key, seed (if account is created on the website) and ot her stuff which
// is yet to be decided

// All investors will be referenced by their public key, name is optional (maybe necessary?)
// we need to stil ldecide on identity and stuff and how much we want to track
// people who invest in the schools
type Investor struct {
	VotingBalance int // this will be equal to the amount of stablecoins that the
	// investor possesses, should update this every once in a while to ensure voting
	// consistency.
	AmountInvested float64
	// total amount, would be nice to track to contact them,
	// give them some kind of medals or something
	InvestedAssets []DBParams
	InvestedBonds []string
	InvestedCoops []string
	// array of asset codes this user has invested in
	// also I think we need a username + password for logging on to the platform itself
	// linking it here for now
	U User
	// user related functions are called as an instance directly
}

// NewInvestor creates a new investor object when passed the username, password hash,
// name and an option to generate the seed and publicKey. This is done because if
// we decide to allow anonymous investors to invest on our platform, we can easily
// insert their publickey into the system and then have hanlders for them signing
// transactions
// TODO: add anonymous investor signing handlers
func NewInvestor(uname string, pwd string, seedpwd string, Name string) (Investor, error) {
	var a Investor
	var err error
	a.U, err = NewUser(uname, pwd, seedpwd, Name)
	if err != nil {
		return a, err
	}
	a.AmountInvested = float64(0)
	err = a.Save()
	return a, err
}

// InsertInvestor inserts a passed Investor object into the database
func (a *Investor) Save() error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(InvestorBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.ItoB(a.U.Index)), encoded)
		// but why do we index based on Index?
		// this is because we do want to enumerate through all investors, which can not be done
		// in a name based construction. But this makes search harder, since now you
		// all entries to find something as simple as a password.
		// TODO: discuss indexing by pwd hash and implications. For small no of entries,
		// we can s till tierate over all the entries.
	})
	return err
}

// RetrieveInvestor retrieves a particular investor indexed by key from the database
func RetrieveInvestor(key int) (Investor, error) {
	var inv Investor
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(InvestorBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return nil
		}
		return json.Unmarshal(x, &inv)
	})
	return inv, err
}

// RetrieveAllInvestors gets a list of all investors in the database
func RetrieveAllInvestors() ([]Investor, error) {
	// this route is broken because it reads through keys sequentially
	// need to see keys until the length of the users database
	var arr []Investor
	temp, err := RetrieveAllUsers()
	if err != nil {
		return arr, err
	}
	limit := len(temp) + 1
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(InvestorBucket)
		for i := 1; i < limit; i++ {
			var rInvestor Investor
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// this is where the key does not exist
				continue
			}
			err := json.Unmarshal(x, &rInvestor)
			if err != nil {
				// we've reached the end of input, so this is not an error
				// ideal error would be "unexpected JSON input" or something similar
				return nil
			}
			arr = append(arr, rInvestor)
		}
		return nil
	})
	return arr, err
}

func ValidateInvestor(name string, pwhash string) (Investor, error) {
	var rec Investor
	user, err := ValidateUser(name, pwhash)
	if err != nil {
		return rec, err
	}
	return RetrieveInvestor(user.Index)
}

func (a *Investor) DeductVotingBalance(votes int) error {
	// TODO: we need to update the voting balance often in accordance with the stablecoin
	// balance or a user will have way less votes. This needs an aadditional field
	// in the db to track past balance and then adjust the amoutn of votes he has
	// accordingly
	a.VotingBalance -= votes
	return a.Save()
}

func (a *Investor) AddVotingBalance(votes int) error {
	// this function is caled when we want to refund the user with the votes once
	// an order has been finalized.
	// TODO: use this
	a.VotingBalance += votes
	return a.Save()
}

// TrustAsset creates a trustline from the caller towards the specific asset
// and asset issuer with a _limit_ set on the maximum amount of tokens that can be sent
// through the trust channel. Each trustline costs 0.5XLM.
func (a *Investor) TrustAsset(asset build.Asset, limit string, seed string) (string, error) {
	// TRUST is FROM recipient TO issuer
	trustTx, err := build.Transaction(
		build.SourceAccount{a.U.PublicKey},
		build.AutoSequence{SequenceProvider: xlm.TestNetClient},
		build.TestNetwork,
		build.Trust(asset.Code, asset.Issuer, build.Limit(limit)),
	)

	if err != nil {
		return "", err
	}

	_, hash, err := xlm.SendTx(seed, trustTx)
	return hash, err
}

// CanInvest checks whether an investor has the required balance to invest in a project
func (a *Investor) CanInvest(balance string, targetBalance string) bool {
	balance, err := xlm.GetUSDTokenBalance(a.U.PublicKey)
	if err != nil {
		return false
	}
	return balance >= targetBalance
}

func (a *Investor) VoteTowardsProposedProject(allProposedProjects []Project, vote int) error {
	// split the coting stuff into a separate function
	// we need to go through the contractor's proposed projects to find an project
	// with index pProjectN
	for _, elem := range allProposedProjects {
		if elem.Params.Index == vote {
			// we have the specific contract and need to upgrade the number of votes on this one
			fmt.Println("YOUR AVAILABLE VOTING BALANCE IS: ", a.VotingBalance)
			fmt.Println("HOW MANY VOTES DO YOU WANT TO DELEGATE TOWARDS THIS ORDER?")
			votes, err := scan.ScanForInt()
			if err != nil {
				return err
			}
			if votes > a.VotingBalance {
				return fmt.Errorf("Can't vote with an amount greater than available balance")
			}
			elem.Params.Votes += votes
			err = elem.Save()
			if err != nil {
				return err
			}
			err = a.DeductVotingBalance(votes)
			if err != nil {
				return err
			}
			fmt.Println("CAST VOTE TOWARDS CONTRACT SUCCESSFULLY")
			log.Println("FOUND CONTRACTOR!")
			return nil
		}
	}
	return fmt.Errorf("Index of project not found, returning")
}

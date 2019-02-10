package opensolar

import (
	"encoding/json"
	"fmt"

	database "github.com/YaleOpenLab/openx/database"
	utils "github.com/YaleOpenLab/openx/utils"
	"github.com/boltdb/bolt"
)

// the contractor super struct comprises of various entities within it. Its a
// super class because combining them results in less duplication of code
type Entity struct {
	// User defines common params such as name, seed, publickey
	U database.User
	// the name of the contractor / company that is contracting
	// A contractor is party who proposes a specific some of money towards a
	// particular project. This is the actual amount that the investors invest in.
	// This ideally must include the developer fee within it, so that investors
	// don't have to invest in two things. It would also make sense because the contractors
	// sometimes would hire developers themselves.
	Contractor bool
	// A guarantor is somebody who can assure investors that the school will get paid
	// on time. This authority should be trusted and either should be vetted by the law
	// or have a multisig paying out to the investors beyond a certain timeline if they
	// don't get paid by the school. This way, the guarantor can be anonymous, like the
	// nice Pineapple Fund guy. This can also be an insurance company, who is willing to
	// guarantee for specific school and the school can pay him out of chain / have
	// that as fee within the contract the originator
	Developer bool
	// A developer is someone who installs the required equipment (Raspberry Pi,
	// network adapters, anti tamper installations and similar) In the initial
	// projects, this will be us, since we'd be installign the pi ourselves, but in
	// the future, we expect third party developers / companies to do this for us
	// and act in a decentralized fashion. This money can either be paid out of chain
	// in fiat or can be a portion of the funds the investors chooses to invest in.
	// a contractor may also employ developers by himself, so this entity is not
	// strictly necessary.
	Originator bool
	// An Originator is an entity that will start a project and get a fixed fee for
	// rendering its service. An Originator's role is not restricted, the originator
	// can also be the developer, contractor or guarantor. The originator should take
	// the responsibility of auditing the requirements of the project - panel size,
	// location, number of panels needed, etc. He then should ideally be able to fill
	// out some kind of form on the website so that the originator's proposal is live
	// and shown to potential investors. The originators get paid only when the project
	// is live, else they can just spam, without any actual investment
	Guarantor bool
	// A Guarantor is someone who can vouch for the recipient and fill in for them
	// in case they default on payment. They can c harge a fee and this must be
	// put inside the contract itself.
	PastContracts []Project
	// list of all the contracts that the contractor has won in the past
	ProposedContracts []Project
	// the Originator proposes a contract which will then be taken up
	// by a contractor, who publishes his own copy of the proposed contract
	// which will be the set of contracts that will be sent to auction
	PresentContracts []Project
	// list of all contracts that the contractor is presently undertaking1
	PastFeedback []Feedback
	// feedback received on the contractor from parties involved in the past
	// What kind of proof do we want from the company? KYC?
	// maybe we could have a photo op like exchanges do these days, with the owner
	// holding up his drivers' license or similar
	Collateral float64
	// the amount of collateral that the entity is willing to hold in case it reneges
	// on a specific contract's details. This is an optional parameter but having collateral
	// would increase investor confidence that a particular entity will keep its word
	// regarding a particular contract.
	CollateralData []string
	// the specific thing(s) which the contractor wants to hold as collateral described
	// as a string (for eg, if a cash bond worht 5000 USD is held as collaterlal,
	// collateral would be set to 5000 USD and CollateralData would be "Cash Bond")
}

func (a *Entity) Save() error {
	db, err := database.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ContractorBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			return err
		}
		return b.Put([]byte(utils.ItoB(a.U.Index)), encoded)
	})
	return err
}

// gets all the proposed contracts for a particular recipient
func RetrieveAllEntitiesWithoutRole() ([]Entity, error) {
	var arr []Entity
	temp, err := database.RetrieveAllUsers()
	if err != nil {
		return arr, err
	}
	limit := len(temp) + 1
	db, err := database.OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ContractorBucket)
		for i := 1; i < limit; i++ {
			var rContractor Entity
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// might be some other user like an investor or recipient
				continue
			}
			err := json.Unmarshal(x, &rContractor)
			if err != nil {
				return nil
			}
			arr = append(arr, rContractor)
		}
		return nil
	})
	return arr, err
}

// gets all the proposed contracts for a particular recipient
func RetrieveAllEntities(role string) ([]Entity, error) {
	var arr []Entity
	temp, err := database.RetrieveAllUsers()
	if err != nil {
		return arr, err
	}
	limit := len(temp) + 1
	db, err := database.OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ContractorBucket)
		for i := 1; i < limit; i++ {
			var rContractor Entity
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// might be some other user like an investor or recipient
				continue
			}
			err := json.Unmarshal(x, &rContractor)
			if err != nil {
				return nil
			}
			switch role {
			case "contractor":
				if !rContractor.Contractor {
					continue
				}
			case "developer":
				if !rContractor.Developer {
					continue
				}
			case "originator":
				if !rContractor.Originator {
					continue
				}
			case "guarantor":
				if !rContractor.Guarantor {
					continue
				}
			default:
				continue
				// default is to add all contractentities to the array
			}
			arr = append(arr, rContractor)
		}
		return nil
	})
	return arr, err
}

func RetrieveEntity(key int) (Entity, error) {
	var a Entity
	db, err := database.OpenDB()
	if err != nil {
		return a, err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ContractorBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return fmt.Errorf("Retrieving entity returns nil, quitting!")
		}
		return json.Unmarshal(x, &a)
	})
	return a, err
}

func newEntity(uname string, pwd string, seedpwd string, Name string, Address string, Description string, role string) (Entity, error) {
	var a Entity
	var err error
	a.U, err = database.NewUser(uname, pwd, seedpwd, Name)
	if err != nil {
		return a, err
	}
	// set all auto fields above
	a.U.Address = Address
	a.U.Description = Description
	// insertion into the database will be a separate handler, pass this Entity there
	switch role {
	case "contractor":
		a.Contractor = true
	case "developer":
		a.Developer = true
	case "originator":
		a.Originator = true
	case "guarantor":
		a.Guarantor = true
	default:
		return a, fmt.Errorf("invalid entity type passed!")
	}
	err = a.Save()
	return a, err
}

func ChangeReputation(entityIndex int, reputation float64) error {
	a, err := RetrieveEntity(entityIndex)
	if err != nil {
		return err
	}
	if reputation > 0 {
		err = a.U.IncreaseReputation(reputation)
	} else {
		err = a.U.DecreaseReputation(reputation)
	}
	if err != nil {
		return err
	}
	return a.Save()
}

func TopReputationEntitiesWithoutRole() ([]Entity, error) {
	// TopReputationEntities returns entities with reputation in descending order
	allEntities, err := RetrieveAllEntitiesWithoutRole()
	if err != nil {
		return allEntities, err
	}
	for i, _ := range allEntities {
		for j, _ := range allEntities {
			if allEntities[i].U.Reputation < allEntities[j].U.Reputation {
				tmp := allEntities[i]
				allEntities[i] = allEntities[j]
				allEntities[j] = tmp
			}
		}
	}
	return allEntities, nil
}

func TopReputationEntities(role string) ([]Entity, error) {
	// caller knows what role he needs this list for, so directly retrieve and do stuff here
	allEntities, err := RetrieveAllEntities(role)
	if err != nil {
		return allEntities, err
	}
	for i, _ := range allEntities {
		for j, _ := range allEntities {
			if allEntities[i].U.Reputation < allEntities[j].U.Reputation {
				tmp := allEntities[i]
				allEntities[i] = allEntities[j]
				allEntities[j] = tmp
			}
		}
	}
	return allEntities, nil
}

func ValidateEntity(name string, pwhash string) (Entity, error) {
	var rec Entity
	user, err := database.ValidateUser(name, pwhash)
	if err != nil {
		return rec, err
	}
	return RetrieveEntity(user.Index)
}

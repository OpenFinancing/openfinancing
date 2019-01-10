package database

import (
	"encoding/json"
	"fmt"
	"log"

	aes "github.com/OpenFinancing/openfinancing/aes"
	utils "github.com/OpenFinancing/openfinancing/utils"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
	"github.com/boltdb/bolt"
)

// the user structure houses all entities that are of type "User". This contains
// commonly used functions so that we need not repeat the ssame thing for every instance.
type User struct {
	Index int
	// default index, gets us easy stats on how many people are there and stuff,
	// don't want to omit this
	EncryptedSeed []byte
	// EncryptedSeed stores the AES-256 encrypted seed of the user. This way, even
	// if the platform is hacked, the user's funds are still safe
	Name string
	// Name of the primary stakeholder involved (principal trustee of school, for eg.)
	PublicKey string
	// PublicKey denotes the public key of the recipient
	LoginUserName string
	// the username you use to login to the platform
	LoginPassword string
	// the password, which you use to authenticate on the platform
	Address string
	// the registered address of the above company
	Description string
	// Does the contractor need to have a seed and a publickey?
	// we assume that it does in this case and proceed.
	// information on company credentials, their experience
	Image string
	// image can be company logo, founder selfie
	// hash of the password in reality
	FirstSignedUp string
	// auto generated timestamp
}

// User is a metastrucutre that contains commonyl used keys within a single umbrella
// so that we can import it wherever needed.
func NewUser(uname string, pwd string, seedpwd string, Name string) (User, error) {
	// call this after the user has failled in username and password.
	// Store hashed password in the database
	var a User

	allUsers, err := RetrieveAllUsers()
	if err != nil {
		return a, err
	}

	// the ugly indexing thing again, need to think of something better here
	if len(allUsers) == 0 {
		a.Index = 1
	} else {
		a.Index = len(allUsers) + 1
	}

	a.Name = Name
	err = a.GenKeys(seedpwd)
	if err != nil {
		return a, err
	}
	a.LoginUserName = uname
	a.LoginPassword = utils.SHA3hash(pwd) // store tha sha3 hash
	// now we have a new User, take this and then send this struct off to be stored in the database
	a.FirstSignedUp = utils.Timestamp()
	err = a.Save()
	return a, err // since user is a meta structure, insert it and then return the function
}

// InsertUser inserts a passed User object into the database
func (a *User) Save() error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.ItoB(a.Index)), encoded)
	})
	return err
}

// RetrieveAllUsers gets a list of all User in the database
func RetrieveAllUsers() ([]User, error) {
	var arr []User
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		for i := 1; ; i++ {
			var rUser User
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rUser)
			if err != nil {
				return err
			}
			arr = append(arr, rUser)
		}
		return nil
	})
	return arr, err
}

// RetrieveUser retrieves a particular User indexed by key from the database
func RetrieveUser(key int) (User, error) {
	var inv User
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return nil
		}
		return json.Unmarshal(x, &inv)
	})
	return inv, err
}

func ValidateUser(name string, pwhash string) (User, error) {
	var inv User
	temp, err := RetrieveAllUsers()
	if err != nil {
		return inv, err
	}
	limit := len(temp) + 1
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		for i := 1; i < limit; i++ {
			var rUser User
			x := b.Get(utils.ItoB(i))
			if x == nil {
				continue
			}
			err := json.Unmarshal(x, &rUser)
			if err != nil {
				return err
			}
			// we have the User class, check names
			if rUser.LoginUserName == name && rUser.LoginPassword == pwhash {
				inv = rUser
				return nil
			}
		}
		return fmt.Errorf("Not Found")
	})
	return inv, err
}

func (a *User) GenKeys(seedpwd string) error {
	var err error
	var seed string
	seed, a.PublicKey, err = xlm.GetKeyPair()
	if err != nil {
		return err
	}
	// don't store the seed in the database
	a.EncryptedSeed, err = aes.Encrypt([]byte(seed), seedpwd)
	err = a.Save()
	return err
}

func (a *User) GetSeed(seedpwd string) (string, error) {
	return wallet.DecryptSeed(a.EncryptedSeed, seedpwd)
}

package database

import (
	"encoding/base32"
	"github.com/pkg/errors"
	"log"

	aes "github.com/Varunram/essentials/aes"
	algorand "github.com/Varunram/essentials/crypto/algorand"
	ethereum "github.com/Varunram/essentials/crypto/eth"
	xlm "github.com/Varunram/essentials/crypto/xlm"
	assets "github.com/Varunram/essentials/crypto/xlm/assets"
	wallet "github.com/Varunram/essentials/crypto/xlm/wallet"
	edb "github.com/Varunram/essentials/database"
	googauth "github.com/Varunram/essentials/googauth"
	recovery "github.com/bithyve/research/sss"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
)

// User is a metastrucutre that contains commonly used keys within a single umbrella
// so that we can import it wherever needed.
type User struct {
	Index int
	// default index, gets us easy stats on how many people are there
	Name string
	// Name of the primary stakeholder involved (principal trustee of school, for eg.)
	StellarWallet  StellWallet
	AlgorandWallet algorand.AlgorandWallet
	// PublicKey denotes the public key of the recipient
	City string
	// the city of residence of the resident
	ZipCode string
	// the zipcode of hte particular city
	Country string
	// the coutnry of residence of the resident
	RecoveryPhone string
	// the phone number where we need to send recovery codes to
	Username string
	// the username you use to login to the platform
	Pwhash string
	// the password hash, which you use to authenticate on the platform
	Address string
	// the registered address of the above company
	Description string
	// Does the contractor need to have a seed and a publickey?
	// we assume that it does in this case and proceed.
	// information on company credentials, their experience
	Image string
	// image can be company logo, founder selfie
	FirstSignedUp string
	// auto generated timestamp
	Kyc bool
	// false if kyc is not accepted / reviewed, true if user has been verified.
	Inspector bool
	// inspector is a kyc inspector who valdiates the data of people who would like
	// to signup on the platform
	Banned bool
	// a field which can be used to set a ban on a user. Can be only used by inspectors in the event someone
	// who has KYC is known to behave in a suspicious way.
	Email string
	// user email to send out notifications
	Notification bool
	// GDPR, if user wants to opt in, set this to true. Default is false
	Reputation float64
	// Reputation contains the max reputation that can be gained by a user. Reputation increases
	// for each completed bond and decreases for each bond cancelled. The frontend
	// could have a table based on reputation scores and use the appropriate scores for
	// awarding badges or something to users with high reputation
	LocalAssets []string
	// a collection of assets that the user can own and trade locally using the emulator
	RecoveryShares []string
	// RecoveryShares are shares that you could hare out to a party and one could reconstruct the
	// seed from 2 out of 3 parts. Based on Shamir's Secret Sharing Scheme.
	PwdResetCode string

	SecondaryWallet Wallet
	// SecondaryWallet defines a higher level wallet which can be imagined to be similar to a savings account

	EthereumWallet ethereum.EthereumWallet
	// EthereumWallet defines a separate wallet for ethereum which people can use to control their ERC721 RECs

	PendingDocuments map[string]string
	// a Pending documents map to keep track of documents that the user in question has to keep track of
	// related to a specific project. The key is the same as the value of the project and the value is a description
	// of what exactly needs to be submitted.
	KYC KycStruct

	StarRating map[int]int // peer bases tarr rating that users can give of each other. Can be gamed, but this is complemented by
	// the automated feedback system, so we should be good.

	GivenStarRating map[int]int // to keep track of users for whom you've given feedback

	TwoFASecret string // the 2FA secret that users can use to authenticate with something like Google Authenticator

	AnchorKYC AnchorKYCHelper // kyc stuff required by AnchorUSD
}

// KycStruct contains the parameters required by the kyc partner for querynig kyc compliance
type KycStruct struct {
	PassportPhoto  string // should be a base64 string or similar according to what the API provider wants
	IDCardPhoto    string
	DriversLicense string
	PersonalPhoto  string // a selfie to verify that  the person registering on the platform is the same person whose documents have been uploaded
}

// AAnchorKYC defines the list of fields that AnchorUSD requires us to provide for getting AnchorUSD
type AnchorKYCHelper struct {
	Name     string
	Birthday struct {
		Month string
		Day   string
		Year  string
	}
	Tax struct {
		Country string
		Id      string
	}
	Address struct {
		Street  string
		City    string
		Postal  string
		Region  string
		Country string
		Phone   string
	}
	PrimaryPhone       string
	Gender             string
	DepositIdentifier  string
	WithdrawIdentifier string
}

type StellWallet struct {
	PublicKey     string
	EncryptedSeed []byte
	SeedPwhash    string
}

// Wallet contains the stuff that we need for a wallet.
type Wallet struct {
	EncryptedSeed []byte // the seedpwd for this would be the same as the one for the primary wallet
	// since we don't want the user to remember like 10 passwords
	PublicKey string
}

// NewUser creates a new user
func NewUser(uname string, pwd string, seedpwd string, Name string) (User, error) {
	// call this after the user has failled in username and password.
	// Store hashed password in the database
	var a User

	allUsers, err := RetrieveAllUsers()
	if err != nil {
		return a, errors.Wrap(err, "Error while retrieving all users from database")
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
		return a, errors.Wrap(err, "Error while generating public and private keys")
	}
	a.Username = uname
	a.Pwhash = utils.SHA3hash(pwd) // store tha sha3 hash
	// now we have a new User, take this and then send this struct off to be stored in the database
	a.FirstSignedUp = utils.Timestamp()
	a.Kyc = false
	a.Notification = false
	err = a.Save()
	return a, err // since user is a meta structure, insert it and then return the function
}

// Save inserts a passed User object into the database
func (a *User) Save() error {
	return edb.Save(consts.DbDir, InvestorBucket, a, a.Index)
}

// RetrieveAllUsersWithoutKyc retrieves all users without kyc
func RetrieveAllUsersWithoutKyc() ([]User, error) {
	var arr []User

	users, err := RetrieveAllUsers()
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, user := range users {
		if !user.Kyc {
			arr = append(arr, user)
		}
	}

	return arr, nil
}

// RetrieveAllUsersWithKyc retrieves all users with kyc
func RetrieveAllUsersWithKyc() ([]User, error) {
	// RetrieveAllUsersWithoutKyc retrieves all users without kyc
	var arr []User

	users, err := RetrieveAllUsers()
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, user := range users {
		if user.Kyc {
			arr = append(arr, user)
		}
	}

	return arr, nil
}

// RetrieveAllUsers gets a list of all User in the database
func RetrieveAllUsers() ([]User, error) {
	var users []User
	x, err := edb.RetrieveAllKeys(consts.DbDir, UserBucket)
	if err != nil {
		return users, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, value := range x {
		users = append(users, value.(User))
	}

	return users, nil
}

// RetrieveUser retrieves a particular User indexed by key from the database
func RetrieveUser(key int) (User, error) {

	var user User
	x, err := edb.Retrieve(consts.DbDir, UserBucket, key)
	if err != nil {
		return user, errors.Wrap(err, "error while retrieving key from bucket")
	}

	return x.(User), nil
}

// ValidateSeedpwd acts as a pre verify function so we don't try to decrypt the encrypted seed
// each time a malicious entity tries to guess the password.
func ValidateSeedpwd(name string, pwhash string, seedpwd string) (User, error) {
	user, err := ValidateUser(name, pwhash)
	if err == nil && utils.SHA3hash(seedpwd) == user.StellarWallet.SeedPwhash {
		return user, nil
	} else {
		return user, errors.New("errored out in seedpwd validation, quitting")
	}
}

// ValidateUser validates a particular user
func ValidateUser(name string, pwhash string) (User, error) {
	var dummy User
	users, err := RetrieveAllUsers()
	if err != nil {
		return dummy, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, user := range users {
		if user.Username == name && user.Pwhash == pwhash {
			return user, nil
		}
	}
	return dummy, errors.New("could not find user with requested credentials")
}

// GenKeys generates a keypair for the user
func (a *User) GenKeys(seedpwd string, options ...string) error {
	if len(options) == 1 {
		chain := options[0]
		switch chain {
		case "algorand":
			log.Println("Generating Algorand wallet")

			var err error
			password := seedpwd

			a.AlgorandWallet, err = algorand.GenNewWallet("algowl", password)
			if err != nil {
				return errors.Wrap(err, "couldn't create new wallet id, quitting")
			}

			err = a.Save()
			if err != nil {
				return err
			}

			backupPhrase, err := algorand.GenerateBackup(a.AlgorandWallet.WalletName, password)
			if err != nil {
				return err
			}

			tmp, err := recovery.Create(2, 3, backupPhrase)
			if err != nil {
				return errors.Wrap(err, "error while storing recovery shares")
			}

			a.RecoveryShares = append(a.RecoveryShares, tmp...) // this is for the primary account
		default:
			log.Println("Chain not supported, please feel free to add support in aanew Pull Request")
			return errors.New("chain not supported, returning")
		} // end of switch
	} else if len(options) == 0 {
		// default user account supported is stellar
		var err error
		var seed string
		seed, a.StellarWallet.PublicKey, err = xlm.GetKeyPair()
		if err != nil {
			return errors.Wrap(err, "error while generating public and private key pair")
		}
		// don't store the seed in the database
		a.StellarWallet.EncryptedSeed, err = aes.Encrypt([]byte(seed), seedpwd)
		if err != nil {
			return errors.Wrap(err, "error while encrypting seed")
		}

		tmp, err := recovery.Create(2, 3, seed)
		if err != nil {
			return errors.Wrap(err, "error while storing recovery shares")
		}

		a.RecoveryShares = append(a.RecoveryShares, tmp...) // this is for the primary account
	}

	secSeed, secPubkey, err := xlm.GetKeyPair()
	if err != nil {
		return errors.Wrap(err, "could not generate secondary keypair")
	}

	a.SecondaryWallet.PublicKey = secPubkey
	a.SecondaryWallet.EncryptedSeed, err = aes.Encrypt([]byte(secSeed), seedpwd)
	if err != nil {
		return errors.Wrap(err, "error while encrypting seed")
	}

	a.EthereumWallet, err = ethereum.GenEthWallet()
	if err != nil {
		return errors.Wrap(err, "error while generating ethereum wallet, quitting")
	}

	err = a.Save()
	return err
}

// CheckUsernameCollision checks if a username is available to a new user who
// wants to signup on the platform
func CheckUsernameCollision(uname string) (User, error) {
	var dummy User
	users, err := RetrieveAllUsers()
	if err != nil {
		return dummy, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, user := range users {
		if user.Username == uname {
			return user, errors.New("username collision observed, quitting")
		}
	}

	return dummy, nil
}

// Authorize authorizes a user
func (a *User) Authorize(userIndex int) error {
	// we don't really mind who this user is since all we need to verify is his identity
	if !a.Inspector {
		return errors.New("You don't have the required permissions to kyc a person")
	}
	user, err := RetrieveUser(userIndex)
	// we want to retrieve only users who have not gone through KYC before
	if err != nil {
		return errors.Wrap(err, "error while retrieving user from database")
	}
	if user.Kyc {
		return errors.New("user already KYC'd")
	}
	user.Kyc = true
	return user.Save()
}

// AddInspector adds a kyc inspector
func AddInspector(userIndex int) error {
	// this should only be called by the platform itself and not open to others
	user, err := RetrieveUser(userIndex)
	if err != nil {
		return errors.Wrap(err, "error while retrieving user from database")
	}
	user.Inspector = true
	return user.Save()
}

// these two functions can be used as internal hnadlers and hte RPC can save reputation directly

// IncreaseReputation increases reputation
func (a *User) ChangeReputation(reputation float64) error {
	a.Reputation += reputation
	return a.Save()
}

// TopReputationUsers gets the users with top reputation
func TopReputationUsers() ([]User, error) {
	// these reputation functions should mostly be used by the frontend through the
	// RPC to display to other users what other users' reputation is.
	allUsers, err := RetrieveAllUsers()
	if err != nil {
		return allUsers, errors.Wrap(err, "error while retrieving all users from database")
	}
	for i := range allUsers {
		for j := range allUsers {
			if allUsers[i].Reputation > allUsers[j].Reputation {
				tmp := allUsers[i]
				allUsers[i] = allUsers[j]
				allUsers[j] = tmp
			}
		}
	}
	return allUsers, nil
}

// IncreaseTrustLimit increases the trustl imit of a specific user towards the STABLEUSD asset
func (a *User) IncreaseTrustLimit(seedpwd string, trust string) error {

	seed, err := wallet.DecryptSeed(a.StellarWallet.EncryptedSeed, seedpwd)
	if err != nil {
		return errors.Wrap(err, "couldn't decrypt seed, quitting!")
	}

	// we now have the seed, so we should upgrade the trustlimit by the margin requested. The margin passed here
	// must not include the old trustlimit

	trustFloat, err := utils.ToFloat(trust)
	if err != nil {
		return err
	}

	stlFloat, err := utils.ToFloat(consts.StablecoinTrustLimit)
	if err != nil {

		return err
	}

	trustLimit, err := utils.ToString(trustFloat + stlFloat)
	if err != nil {
		return err
	}
	_, err = assets.TrustAsset(consts.StablecoinCode, consts.StableCoinAddress, trustLimit, seed)
	if err != nil {
		return errors.Wrap(err, "couldn't trust asset, quitting!")
	}

	return nil
}

// SearchWithEmailId searches for a given user who has the given email id
func SearchWithEmailId(email string) (User, error) {
	var dummy User
	users, err := RetrieveAllUsers()
	if err != nil {
		return dummy, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}

	return dummy, errors.New("could not find user with requested email id, quitting")
}

// MoveFundsFromSecondaryWallet moves funds from the secondary wallet to the primary wallet
func (a *User) MoveFundsFromSecondaryWallet(amount string, seedpwd string) error {
	amountI, err := utils.ToFloat(amount)
	if err != nil {
		return errors.Wrap(err, "amount not float, quitting")
	}
	// unlock secondary account
	secSeed, err := wallet.DecryptSeed(a.SecondaryWallet.EncryptedSeed, seedpwd)
	if err != nil {
		return errors.Wrap(err, "could not unlock secondary seed, quitting")
	}

	// get secondary balance
	secFunds, err := xlm.GetNativeBalance(a.SecondaryWallet.PublicKey)
	if err != nil {
		return errors.Wrap(err, "could not get xlm balance of secondary account")
	}

	secFundsFloat, err := utils.ToFloat(secFunds)
	if err != nil {
		return err
	}

	if amountI > secFundsFloat {
		return errors.New("amount to be transferred is greater than the funds available in the secondary account, quitting")
	}

	// send the tx over
	_, txhash, err := xlm.SendXLM(a.StellarWallet.PublicKey, amount, secSeed, "fund transfer to secondary")
	if err != nil {
		return errors.Wrap(err, "error while transferring funds to secondary account, quitting")
	}

	log.Println("transfer sec-prim tx hash: ", txhash)
	return nil
}

// SweepSecondaryWallet sweeps fudsd from the secondary account to the primary account
func (a *User) SweepSecondaryWallet(seedpwd string) error {
	// unlock secondary account

	secSeed, err := wallet.DecryptSeed(a.SecondaryWallet.EncryptedSeed, seedpwd)
	if err != nil {
		return errors.Wrap(err, "could not unlock primary seed, quitting")
	}

	// get secondary balance
	secFunds, err := xlm.GetNativeBalance(a.SecondaryWallet.PublicKey)
	if err != nil {
		return errors.Wrap(err, "could not get xlm balance of secondary account")
	}

	secFundsTemp, err := utils.ToFloat(secFunds)
	if err != nil {
		return err
	}

	secFundsWithMinbal, err := utils.ToString(secFundsTemp - 5)
	if err != nil {
		return err
	}
	// send the tx over
	_, txhash, err := xlm.SendXLM(a.StellarWallet.PublicKey, secFundsWithMinbal, secSeed, "fund transfer to secondary")
	if err != nil {
		return errors.Wrap(err, "error while transferring funds to secondary account, quitting")
	}

	log.Println("transfer sec-prim tx hash: ", txhash)
	return nil
}

// AddEmail stores the passed email as the user's email.
func (a *User) AddEmail(email string) error {
	// call this function when a user wants to get notifications. Ask on frontend whether
	// it wants to
	a.Email = email
	a.Notification = true
	err := a.Save()
	if err != nil {
		return errors.Wrap(err, "error while saving investor")
	}
	return a.Save()
}

// SetBan can be used by an inspector to se a ban on any user for violating certain terms and conditions
func (a *User) SetBan(userIndex int) error {
	if !a.Inspector {
		return errors.New("user not authorized to ban a user")
	}

	if a.Index == userIndex {
		return errors.New("can't ban yourself, quitting")
	}

	user, err := RetrieveUser(userIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't  find user to ban, quitting")
	}

	if user.Banned {
		return errors.Wrap(err, "user already banned, not setitng another ban")
	}

	user.Banned = true
	return user.Save()
}

// GiveFeedback is used by a user to give a star based feedback about the other user
func (a *User) GiveFeedback(userIndex int, feedback int) error {
	user, err := RetrieveUser(userIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve user from db while giving feedback")
	}

	if len(user.StarRating) == 0 {
		// no one has given t3his user a starr rating before, so create a new map
		user.StarRating = make(map[int]int)
	}

	if feedback > 5 || feedback < 0 {
		log.Println("feedback greater than 5 or less than 0, quitting")
		return errors.New("feedback greater than 5, quitting")
	}

	user.StarRating[a.Index] = feedback
	log.Println("STARRATING: ", user.StarRating, user.Name)
	err = user.Save()
	if err != nil {
		return errors.Wrap(err, "couldn't save feedback provided on user")
	}

	if len(a.GivenStarRating) == 0 {
		// no one has given t3his user a starr rating before, so create a new map
		a.GivenStarRating = make(map[int]int)
	}

	a.GivenStarRating[user.Index] = feedback
	return a.Save()
}

// Generate2FA generates a new 2FA secret for the given user
func (a *User) Generate2FA() (string, error) {
	secret := utils.GetRandomString(35) // multiples of 5 to  prevent the = padding at the end
	secretBase32 := base32.StdEncoding.EncodeToString([]byte(secret))
	otpc := &googauth.OTPConfig{
		Secret:     secretBase32,
		WindowSize: 1,
		UTC:        true,
	}
	otpString, err := otpc.GenerateURI(a.Name)
	if err != nil {
		return otpString, err
	}
	if err != nil {
		return otpString, err
	}
	a.TwoFASecret = secret
	err = a.Save()
	if err != nil {
		return otpString, err
	}
	return otpString, nil
}

// Authenticate2FA authenticates the given password against the user's stored password
func (a *User) Authenticate2FA(password string) (bool, error) {
	secretBase32 := base32.StdEncoding.EncodeToString([]byte(a.TwoFASecret))
	otpc := &googauth.OTPConfig{
		Secret:     secretBase32,
		WindowSize: 1,
		UTC:        true,
	}

	return otpc.Authenticate(password)
}

// +build all travis

package database

import (
	"log"
	"os"
	"testing"

	consts "github.com/OpenFinancing/openfinancing/consts"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
	"github.com/stellar/go/build"
)

// TODO: rewrite how this works and split between platforms and database
// go test --tags="all" -coverprofile=test.txt .
func TestDb(t *testing.T) {
	var err error
	CreateHomeDir()                     // create home directory if it doesn't exist yet
	os.Remove(consts.DbDir + "/yol.db") // remove the database file, if it exists
	consts.DbDir = "blah"               // set to a false db so that we can test errors arising from OpenDB()
	_, err = OpenDB()
	if err == nil { // wrong dir, so should error out
		t.Fatalf("Able to open database with wrong path")
	}
	// The following tests should fail because the db path is invalid and OpenDB() would fail
	err = DeleteKeyFromBucket(1, UserBucket)
	if err == nil {
		t.Fatalf("Able to delete from database with wrong path")
	}
	inv, err := NewInvestor("investor1", "blah", "blah", "Investor1")
	if err == nil {
		t.Fatalf("Able to create investor in database with wrong path")
	}
	var i1 Investor
	err = i1.Save()
	if err == nil {
		t.Fatalf("Able to create investor in database with wrong path")
	}
	_, err = RetrieveInvestor(1)
	if err == nil {
		t.Fatalf("Able to retrieve investor in database with wrong path")
	}
	_, err = RetrieveAllInvestors()
	if err == nil {
		t.Fatalf("Able to retrieve investors in an invalid db, quitting!")
	}
	recp, err := NewRecipient("recipient1", "blah", "blah", "Recipient1")
	if err == nil {
		t.Fatalf("Able to create recipient in database with wrong path")
	}
	var r1 Recipient
	err = r1.Save()
	if err == nil {
		t.Fatalf("Able to save recipient in database with wrong path")
	}
	_, err = RetrieveAllRecipients()
	if err == nil {
		t.Fatalf("Able to retrieve all recipients in database with wrong path")
	}
	_, err = RetrieveRecipient(1)
	if err == nil {
		t.Fatalf("Able to retrieve recipient in database with wrong path")
	}
	user, err := NewUser("user1", "blah", "blah", "User1")
	if err == nil {
		t.Fatalf("Able to retrieve create user in database with wrong path")
	}
	var u1 User
	err = u1.Save()
	if err == nil {
		t.Fatalf("Able to create user in database with wrong path")
	}
	_, err = RetrieveAllUsers()
	if err == nil {
		t.Fatalf("Able to retrieve all users in database with wrong path")
	}
	_, err = RetrieveUser(1)
	if err == nil {
		t.Fatalf("Able to retrieve user in database with wrong path")
	}
	_, err = ValidateUser("blah", "blah")
	if err == nil {
		t.Fatalf("Able to validate user in database with wrong path")
	}
	// set the db directory back to normal so that we can test stuff which goes inside the db
	consts.DbDir = os.Getenv("HOME") + "/.openfinancing/database"
	err = os.MkdirAll(consts.DbDir, os.ModePerm) // create the db
	if err != nil {
		t.Fatal(err)
	}
	// we need to check if we error out while creating buckets. The only way to do that
	// is to set the bucket names to an invalid string so that boltdb errors out and
	// we try to catch that error. Bit ugly, but no other wya than setting and unsetting names
	ProjectsBucket = []byte("")
	x, err := OpenDB()
	if err == nil {
		t.Fatalf("Invalid bucket name")
	}
	x.Close()
	ProjectsBucket = []byte("Projects")
	InvestorBucket = []byte("")
	x, err = OpenDB()
	if err == nil {
		t.Fatalf("Invalid bucket name")
	}
	x.Close()
	InvestorBucket = []byte("Investors")
	RecipientBucket = []byte("")
	x, err = OpenDB()
	if err == nil {
		t.Fatalf("Invalid bucket name")
	}
	x.Close()
	RecipientBucket = []byte("Recipients")
	ContractorBucket = []byte("")
	x, err = OpenDB()
	if err == nil {
		t.Fatalf("Invalid bucket name")
	}
	x.Close()
	ContractorBucket = []byte("Contractors")
	UserBucket = []byte("")
	x, err = OpenDB()
	if err == nil {
		t.Fatalf("Invalid bucket name")
	}
	x.Close()
	UserBucket = []byte("Users")
	BondBucket = []byte("")
	x, err = OpenDB()
	if err == nil {
		t.Fatalf("Invalid bucket name")
	}
	x.Close()
	BondBucket = []byte("Bonds")
	CoopBucket = []byte("")
	x, err = OpenDB()
	if err == nil {
		t.Fatalf("Invalid bucket name")
	}
	x.Close()
	// even though we set the names back to their originals above, ahve this snippet
	// here so that its easier to audit the tests without having to worry about
	// typos while setting the bucket names back to what they were
	CoopBucket = []byte("Coop")
	ProjectsBucket = []byte("Projects")
	InvestorBucket = []byte("Investors")
	RecipientBucket = []byte("Recipients")
	ContractorBucket = []byte("Contractors")
	UserBucket = []byte("Users")
	BondBucket = []byte("Bonds")
	CoopBucket = []byte("Coop")
	db, err := OpenDB()
	if err != nil {
		t.Fatal(err)
	}
	db.Close() // close immmediately after check
	inv, err = NewInvestor("investor1", "blah", "blah", "Investor1")
	if err != nil {
		t.Fatal(err)
	}
	xc, err := RetrieveAllInvestors()
	if len(xc) != 1 {
		t.Fatalf("ERROR!")
	}
	// try retrieving existing stuff
	inv1, err := RetrieveInvestor(1)
	if err != nil {
		t.Fatal(err)
	}
	if inv1.U.Name != "Investor1" {
		t.Fatalf("Investor names don't match, quitting!")
	}
	// func NewRecipient(uname string, pwd string, seedpwd string, Name string) (Recipient, error) {
	recp, err = NewRecipient("recipient1", "blah", "blah", "Recipient1")
	if err != nil {
		t.Fatal(err)
	}

	rec1, err := RetrieveRecipient(recp.U.Index)
	if err != nil {
		t.Fatal(err)
	}

	if rec1.U.Name != "Recipient1" {
		t.Fatalf("Recipient usernames don't match. quitting!")
	}

	user, err = NewUser("user1", "blah", "blah", "User1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = RetrieveInvestor(1000)
	if err == nil {
		t.Fatalf("Investor shouldn't exist, but does, quitting!")
	}

	_, err = RetrieveRecipient(1000)
	if err == nil {
		t.Fatalf("Investor shouldn't exist, but does, quitting!")
	}

	user1, err := RetrieveUser(user.Index)
	if err != nil {
		t.Fatal(err)
	}
	if user1.Name != "User1" {
		t.Fatalf("Usernames don't match. quitting!")
	}

	tmpuser, err := RetrieveUser(1000)
	if tmpuser.Index != 0 {
		t.Fatalf("Investor shouldn't exist, but does, quitting!")
	}

	// test length of users in bucket
	allInv, err := RetrieveAllInvestors()
	if err != nil {
		t.Fatal(err)
	}
	if len(allInv) != 1 {
		log.Println("UNKNOWN: ", len(allInv))
		t.Fatalf("Unknown investors existing, quitting!")
	}

	allRec, err := RetrieveAllRecipients()
	if err != nil {
		t.Fatal(err)
	}
	if len(allRec) != 1 {
		t.Fatalf("Unknown recipients existing, quitting!")
	}

	allUser, err := RetrieveAllUsers()
	if err != nil {
		t.Fatal(err)
	}
	if len(allUser) != 3 {
		t.Fatalf("Unknown users existing, quitting!")
	}

	// check if each of the validate functions work
	_, err = ValidateInvestor("investor1", "ed2df20bb16ecb0b4b149cf8e7d9819afd608b22999e707364196187fca0cf38544c9f3eb981ad81cef18562e4c818370eab068992639af7d70488945265197f")
	if err != nil {
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	_, err = ValidateRecipient("recipient1", "ed2df20bb16ecb0b4b149cf8e7d9819afd608b22999e707364196187fca0cf38544c9f3eb981ad81cef18562e4c818370eab068992639af7d70488945265197f")
	if err != nil {
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	_, err = ValidateUser("user1", "ed2df20bb16ecb0b4b149cf8e7d9819afd608b22999e707364196187fca0cf38544c9f3eb981ad81cef18562e4c818370eab068992639af7d70488945265197f")
	if err != nil {
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	_, err = ValidateInvestor("blah", "ed2df20bb16ecb0b4b149cf8e7d9819afd608b22999e707364196187fca0cf38544c9f3eb981ad81cef18562e4c818370eab068992639af7d70488945265197f")
	if err == nil {
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	_, err = ValidateRecipient("blah", "ed2df20bb16ecb0b4b149cf8e7d9819afd608b22999e707364196187fca0cf38544c9f3eb981ad81cef18562e4c818370eab068992639af7d70488945265197f")
	if err == nil {
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	_, err = ValidateUser("blah", "ed2df20bb16ecb0b4b149cf8e7d9819afd608b22999e707364196187fca0cf38544c9f3eb981ad81cef18562e4c818370eab068992639af7d70488945265197f")
	if err == nil {
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	// check voting balance routes
	voteBalance := inv.VotingBalance
	err = inv.AddVotingBalance(10000)
	if err != nil {
		t.Fatal(err)
	}
	if inv.VotingBalance-voteBalance != 10000 {
		t.Fatalf("Voting Balance not added, quitting!")
	}
	err = inv.DeductVotingBalance(10000)
	if err != nil {
		t.Fatal(err)
	}
	if inv.VotingBalance-voteBalance != 0 {
		t.Fatalf("Voting Balance not added, quitting!")
	}

	// check CanInvest Route
	if inv.CanInvest("100", "1000") {
		t.Fatalf("CanInvest Returns true!")
	}

	err = user.GenKeys("blah")
	if err != nil {
		t.Fatalf("Not able to generate keys, quitting!")
	}

	_, err = user.GetSeed("blah")
	if err != nil {
		t.Fatal(err)
	}

	err = DeleteKeyFromBucket(user.Index, UserBucket)
	if err != nil {
		t.Fatal(err)
	}

	// check the asset functions below. For some weird reason, placing these tests
	// above confuses the other routes, so placing everything here so that we can
	// isolate them from the other routes.
	err = xlm.GetXLM(recp.U.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	err = xlm.GetXLM(inv.U.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	testAsset := build.CreditAsset("blah", recp.U.PublicKey)
	invSeed, err := wallet.DecryptSeed(inv.U.EncryptedSeed, "blah")
	if err != nil {
		t.Fatal(err)
	}
	hash, err := inv.TrustAsset(testAsset, "100", invSeed)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("HASH IS: ", hash)
	_, err = inv.TrustAsset(testAsset, "-1", "blah")
	if err == nil {
		t.Fatalf("can trust asset with invalid s eed!")
	}
	recpSeed, err := wallet.DecryptSeed(recp.U.EncryptedSeed, "blah")
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = recp.SendAssetToIssuer(testAsset.Code, recp.U.PublicKey, "1", recpSeed)
	if err != nil {
		t.Fatal(err)
	}
	pkSeed, pk, err := xlm.GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = recp.SendAssetToIssuer(testAsset.Code, pk, "1", pkSeed) // should fail because
	if err == nil {
		t.Fatalf("Invalid tx succeeds, quitting!")
	}
	testAsset2 := build.CreditAsset("blah2", pkSeed) // this account doesn't exist yet, so this should fail
	_, err = inv.TrustAsset(testAsset2, "-1", "blah")
	if err == nil {
		t.Fatalf("can trust invalid asset")
	}
	os.Remove(consts.DbDir + "/yol.db")
}

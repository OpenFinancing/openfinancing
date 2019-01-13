package main

import (
	aes "github.com/OpenFinancing/openfinancing/aes"
	ipfs "github.com/OpenFinancing/openfinancing/ipfs"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
	"log"
)

type Client struct {
	Info     string
	Location string
	UniqueId string
}

// we should Authenticate whether the code we wanted the rpi to run is actually running
// in order to do this, we have a password + pk/sk based authentication system
var StartHash string
var NowHash string

// Authenticate authenticates a teller. If pubkey doesn't match with hardcoded pk,
// do not start since device might be tampered
func Authenticate() {
	data, err := aes.DecryptFile("auth/auth.txt", "password")
	if err != nil {
		log.Fatal(err)
	}
	pubkey, err := wallet.ReturnPubkey(string(data))
	if err != nil {
		log.Fatal(err)
	}
	if pubkey != "GAJJFLMAPBIORAA3CGS5PKSVB4KG52HX2SMEA2AZWIHNGHW5PWHXLXF7" {
		// this part must be changed each time we are installing a new unit
		log.Fatalf("Unauthorized access, quitting!")
	}
	StartHash, err = BlockStamp()
	if err != nil {
		log.Fatal(err)
	}
	_, err = xlm.GetNativeBalance(pubkey)
	if err != nil {
		// no funds
		err = xlm.GetXLM(pubkey)
		if err != nil {
			log.Fatal(err)
		}
	}
	PublicKey = pubkey
	Seed = string(data)
}

// New generates a new password pair that should be used to authenticate a teller
func New() {
	// call this when you need to generate a keypair for authentication
	seed, _, err := xlm.GetKeyPair()
	if err != nil {
		log.Fatal(err)
	}
	// encrypt and store this seed
	// log.Println("ADDRESS IS: ", address)
	err = aes.EncryptFile("auth/auth.txt", []byte(seed), "password")
	if err != nil {
		log.Fatal(err)
	}
}

// BlockStamp gets the latest blockhash from the API
func BlockStamp() (string, error) {
	// get the latest  block here
	hash, err := xlm.GetLatestBlockHash()
	return hash, err
}

// EndHandler runs when the teller shuts down
func EndHandler(t Client) error {
	log.Println("Gracefully shutting down, please do not press any buttons in the process")
	var err error
	NowHash, err = BlockStamp()
	if err != nil {
		return err
	}
	log.Printf("StartHash: %s, NowHash: %s", StartHash, NowHash)
	hashString := "Device Shutting down. Info: " + t.Info + " Device Location: " + t.Location + " Device Unique ID: " + t.UniqueId + " " + StartHash + NowHash
	// need to hash this with ipfs
	ipfsHash, err := ipfs.AddStringToIpfs(hashString)
	if err != nil {
		return err
	}
	log.Println("ipfs hash: ", ipfsHash)
	memoText := "IPFSHASH: " + ipfsHash
	// 10 + 46 (ipfs hash length) characters
	firstHalf := memoText[:28]
	secondHalf := memoText[28:]
	log.Println("PUBKEY: ", PublicKey)
	_, tx1, err := xlm.SendXLM(PublicKey, "1", Seed, firstHalf)
	if err != nil {
		return err
	}
	_, tx2, err := xlm.SendXLM(PublicKey, "1", Seed, secondHalf)
	if err != nil {
		return err
	}
	log.Printf("tx1 hash: %s, tx2 hash: %s", tx1, tx2)
	return nil
}

// so the teller will be run on the hub and has some data that the platform might need
//  the teller must serve some data to other entities as well. So we need a server for that
// and this must be over tls for preventing mitm attacks and a good tls certificate from an authorized
// provider

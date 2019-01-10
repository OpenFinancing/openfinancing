package database

import (
	utils "github.com/OpenFinancing/openfinancing/utils"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
	"github.com/stellar/go/build"
)

// CreateAsset creates a new asset belonging to the public key referenced above
func CreateAsset(assetName string, PublicKey string) build.Asset {
	// need to set a couple flags here
	return build.CreditAsset(assetName, PublicKey)
}

func AssetID(inputString string) string {
	// so the assetID right now is a hash of the asset name, concatenated investor public keys and nonces
	x := utils.SHA3hash(inputString)
	return "YOL" + x[64:73] // max length of an asset in stellar is 12
	// log.Fatal(fmt.Errorf("All good"))
	// return nil
}

// TrustAsset trusts a specific asset issued by a particular public key and signs
// a transaction with a preset limit on how much it is willing to trsut that issuer's
// asset for
func TrustAsset(asset build.Asset, limit string, PublicKey string, Seed string) (string, error) {
	// TRUST is FROM recipient TO issuer
	trustTx, err := build.Transaction(
		build.SourceAccount{PublicKey},
		build.AutoSequence{SequenceProvider: xlm.TestNetClient},
		build.TestNetwork,
		build.Trust(asset.Code, asset.Issuer, build.Limit(limit)),
	)

	_, txHash, err := xlm.SendTx(Seed, trustTx)
	return txHash, err
}


// SendAsset transfers _amount_ number of assets from the caller to the destination
// and returns an error if the destination doesn't have a trustline with the issuer
// This method is called by the issuer of the asset
func SendAssetFromIssuer(assetName string, destination string, amount string, Seed string, PublicKey string) (int32, string, error) {
	// this transaction is FROM issuer TO recipient
	paymentTx, err := build.Transaction(
		build.SourceAccount{PublicKey},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: xlm.TestNetClient},
		build.MemoText{"Sending Asset: " + assetName},
		build.Payment(
			build.Destination{AddressOrSeed: destination},
			build.CreditAmount{assetName, PublicKey, amount},
			// CreditAmount identifies the asset by asset Code and issuer pubkey
		),
	)

	if err != nil {
		return -1, "", err
	}
	return xlm.SendTx(Seed, paymentTx)
}

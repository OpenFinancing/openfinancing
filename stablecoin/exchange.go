package stablecoin

import (
	assets "github.com/OpenFinancing/openfinancing/assets"
	consts "github.com/OpenFinancing/openfinancing/consts"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
	"log"
)

func Exchange(recipientPK string, recipientSeed string, convAmount string) error {
	hash, err := assets.TrustAsset(StableUSD, consts.StablecoinTrustLimit, recipientPK, recipientSeed)
	if err != nil {
		return err
	}
	log.Println("tx hash for trusting stableUSD: ", hash)
	// now send coins across and see if our tracker detects it
	_, hash, err = xlm.SendXLM(PublicKey, convAmount, recipientSeed, "Sending XLM to bootstrap")
	if err != nil {
		return err
	}
	log.Println("tx hash for sent xlm: ", hash, "pubkey: ", recipientPK)
	return nil
}

package models

import (
	"github.com/pkg/errors"
	"log"
	"time"

	xlm "github.com/YaleOpenLab/openx/chains/xlm"
	assets "github.com/YaleOpenLab/openx/chains/xlm/assets"
	consts "github.com/YaleOpenLab/openx/consts"
)

// the models package won't be imported directly in any place but would be imported
// by all the investment models that exist

// SendUSDToPlatform sends STABLEUSD back to the platform for investment
func SendUSDToPlatform(invSeed string, invAmount float64, memo string) (string, error) {
	// send stableusd to the platform (not the issuer) since the issuer will be locked
	// and we can't use the funds. We also need ot be able to redeem the stablecoin for fiat
	// so we can't burn them
	var oldPlatformBalance float64
	var err error
	oldPlatformBalance, err = xlm.GetAssetBalance(consts.PlatformPublicKey, consts.StablecoinCode)
	if err != nil {
		// platform does not have stablecoin, shouldn't arrive here ideally
		oldPlatformBalance = 0
	}

	var txhash string
	if !consts.Mainnet {
		_, txhash, err = assets.SendAsset(consts.StablecoinCode, consts.StablecoinPublicKey, consts.PlatformPublicKey, invAmount, invSeed, memo)
		if err != nil {
			return txhash, errors.Wrap(err, "sending stableusd to platform failed")
		}
	} else {
		_, txhash, err = assets.SendAsset(consts.AnchorUSDCode, consts.AnchorUSDAddress, consts.PlatformPublicKey, invAmount, invSeed, memo)
		if err != nil {
			return txhash, errors.Wrap(err, "sending stableusd to platform failed")
		}
	}

	log.Println("Sent STABLEUSD to platform, confirmation: ", txhash)
	time.Sleep(5 * time.Second) // wait for a block

	var newPlatformBalance float64
	if !consts.Mainnet {
		newPlatformBalance, err = xlm.GetAssetBalance(consts.PlatformPublicKey, consts.StablecoinCode)
		if err != nil {
			return txhash, errors.Wrap(err, "error while getting asset balance")
		}
	} else {
		newPlatformBalance, err = xlm.GetAssetBalance(consts.PlatformPublicKey, consts.AnchorUSDCode)
		if err != nil {
			return txhash, errors.Wrap(err, "error while getting asset balance")
		}
	}

	if newPlatformBalance-oldPlatformBalance < invAmount-1 {
		return txhash, errors.New("Sent amount doesn't match with investment amount")
	}
	return txhash, nil
}

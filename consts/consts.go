package consts

import (
	"os"
)

var DonateBalance = "10"                                                             // we send this amount free to invesotrs who signup on our platform to enable them to have trustlines. Maybe we should have a payment provider and take money from them?
var StablecoinTrustLimit = "100000"                                                  // the maximum limit that the investor trusts the stablecoin issuer for
var INVAssetPrefix = "BondTokens_"                                                   // the prefix that will be hashed to give an investor AssetID
var BondAssetPrefix = "BondTokens_"                                                  // the prefix that will be hashed to give a bond asset
var CoopAssetPrefix = "CoopAsset_"                                                   // the prefix that will be hashed to give the cooperative asset
var DEBAssetPrefix = "DEBTokens_"                                                    // the prefix that will be hashed to give a recipient AssetID
var PBAssetPrefix = "PBTokens_"                                                      // the prefix that will be hashed to give a payback AssetID
var HomeDir = os.Getenv("HOME") + "/.openfinancing"                                  // home directory where we store the platform seed
var PlatformSeedFile = HomeDir + "/platformseed.hex"                                 // the path where the platform's seed is stored
var StableCoinSeedFile = HomeDir + "/stablecoinseed.hex"                             // the path where the stablecoin's seed is stored
var DbDir = HomeDir + "/database"                                                    // the directory where the main assets of our platform are stored
var IpfsFileLength = 10                                                              // the length of the hash that we want our ipfs hashes to have
const StableCoinAddress = "GBY3DHWSN5CHJ5FDHD7PI5Q23NNMJAK7MGRSERKMOV6QBR7IMAI3IWK5" // the address of the stabellcoin must be a constant for the payment listener to work properly
var TellerHomeDir = os.Getenv("HOME") + "/.openfinancing/teller"                     // the hom directory of the teller executable

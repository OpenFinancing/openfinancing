package main

import (
	"fmt"
	"log"
	"os"

	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	// ipfs "github.com/YaleOpenLab/openx/ipfs"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
	rpc "github.com/YaleOpenLab/openx/rpc"
	// scan "github.com/YaleOpenLab/openx/scan"
	stablecoin "github.com/YaleOpenLab/openx/stablecoin"
	utils "github.com/YaleOpenLab/openx/utils"
	// wallet "github.com/YaleOpenLab/openx/wallet"
	// xlm "github.com/YaleOpenLab/openx/xlm"
	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
)

// the server powering the openx platform of platforms. There are two clients that can be used
// with the backned - ofcli and emulator
// refer https://github.com/stellar/go/blob/master/build/main_test.go in case the stellar
// go SDK docs are insufficient.
var opts struct {
	Port     int  `short:"p" description:"The port on which the server runs on"`
	Simulate bool `short:"t" description:"Simulate the test database with demo values"`
}

// ParseConfig parses CLI parameters passed
func ParseConfig(args []string) (string, error) {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return "", err
	}
	port := utils.ItoS(consts.DefaultRpcPort)
	if opts.Port != 0 {
		port = utils.ItoS(opts.Port)
		log.Println("Starting RPC Server on Port: ", opts.Port)
	}
	return port, nil
}

// StartPlatform starts the platform
func StartPlatform() error {

	database.CreateHomeDir()
	var err error
	// init stablecoin before platform so we don't have to create a stablecoin in case our dbdir is wiped
	err = stablecoin.InitStableCoin() // start the stablecoin daemon
	if err != nil {
		return err
	}

	err = opensolar.InitializePlatform()
	if err != nil {
		return err
	}

	allContracts, err := opensolar.RetrieveAllProjects()
	if err != nil {
		log.Println("Error retrieving all projects from the database")
		return err
	}

	if len(allContracts) == 0 {
		log.Println("Populating database with test values")
		err = InsertDummyData(opts.Simulate)
		if err != nil {
			return err
		}
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Error while reading platform email from config file")
		return err
	}
	consts.PlatformEmail = viper.Get("platformemail").(string)
	log.Println("PLATFORM EMAIL: ", consts.PlatformEmail)
	return nil
}

func main() {
	var err error
	port, err := ParseConfig(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	err = StartPlatform()
	if err != nil {
		log.Fatal(err)
	}

	err = parseYamlProject("100kwy", "data-sandbox/100kw.json", 8)
	if err != nil {
		log.Fatal(err)
	}

	go opensolar.MonitorTeller(1)

	fmt.Printf("PLATFORM SEED IS: %s\n PLATFORM PUBLIC KEY IS: %s\n", consts.PlatformSeed, consts.PlatformPublicKey)
	fmt.Printf("STABLECOIN PUBLICKEY IS: %s\nSTABLECOIN SEED is: %s\n\n", consts.StablecoinPublicKey, consts.StablecoinSeed)
	rpc.StartServer(port)
}

package rpc

import (
	"encoding/json"
	"log"
	"net/http"

	database "github.com/OpenFinancing/openfinancing/database"
	bonds "github.com/OpenFinancing/openfinancing/platforms/bonds"
	utils "github.com/OpenFinancing/openfinancing/utils"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
)

func setupCoopRPCs() {
	getCoopDetails()
	InvestInCoop()
}

func getCoopDetails() {
	http.HandleFunc("/coop/get", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		// get the details of a specific bond by key
		if r.URL.Query()["index"] == nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		uKey := utils.StoI(r.URL.Query()["index"][0])
		bond, err := bonds.RetrieveCoop(uKey)
		if err != nil {
			log.Println(err)
		}
		bondJson, err := json.Marshal(bond)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, bondJson)
	})
}

// curl request attached for convenience
// curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -H "Origin: localhost" -H "Cache-Control: no-cache" -d 'MonthlyPayment=1000&CoopIndex=1&InvIndex=2&InvSeedPwd=x' "http://localhost:8080/coop/invest"
func InvestInCoop() {
	http.HandleFunc("/coop/invest", func(w http.ResponseWriter, r *http.Request) {
		checkPost(w, r)
		var err error
		var iCoop bonds.Coop
		// need to receive a whole lot of parameters here
		// need the bond index passed so that we can retrieve the bond easily
		if r.FormValue("MonthlyPayment") == "" || r.FormValue("CoopIndex") == "" || r.FormValue("InvIndex") == "" || r.FormValue("InvSeedPwd") == "" {
			log.Println("missing params")
			errorHandler(w, r, http.StatusNotFound)
		}

		// TODO: change that this is hardcoded, there must be a nicer way to do this
		// maybe read from the seed and try to decrypt?
		issuerSeed := "SBBYVEI4YNKZANRQEFH35U5GPEJ27MBLL7XHEKX5VC75QLJZWAXGX36Y"
		issuerPk := "GAEY5TVFYWBIIHF7PQCQVNIFTNIF7QSG4IH27HRW3DH476RI4NA2BPV3"
		_, err = xlm.GetNativeBalance(issuerPk)
		if err != nil {
			err = xlm.GetXLM(issuerPk)
			if err != nil {
				log.Fatal(err)
			}
		}

		invAmount := r.FormValue("MonthlyPayment")
		CoopIndex := utils.StoI(r.FormValue("CoopIndex"))
		invIndex := utils.StoI(r.FormValue("InvIndex"))
		invSeedPwd := r.FormValue("InvSeedPwd")

		iCoop, err = bonds.RetrieveCoop(CoopIndex)
		if err != nil {
			log.Fatal(err)
		}
		// pass the investor index, pk and seed
		iInv, err := database.RetrieveInvestor(invIndex)
		if err != nil {
			log.Fatal(err)
		}

		_, err = xlm.GetNativeBalance(iInv.U.PublicKey) // get testnet funds if their account is new
		if err != nil {
			err = xlm.GetXLM(iInv.U.PublicKey)
			if err != nil {
				log.Fatal(err)
			}
		}
		invSeed, err := iInv.U.GetSeed(invSeedPwd)
		if err != nil {
			log.Println("Error while getting seed, inv")
			log.Fatal(err)
		}
		err = iCoop.Invest(issuerPk, issuerSeed, &iInv, invAmount, invSeed)
		if err != nil {
			log.Println(err)
		}
		log.Println("UPDATED BOND: ", iCoop)
		bondJson, err := json.Marshal(iCoop)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, bondJson)
	})
}

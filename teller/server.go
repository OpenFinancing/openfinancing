package main

import (
	"encoding/json"
	//"fmt"
	"log"
	"net/http"

	consts "github.com/YaleOpenLab/openx/consts"
	rpc "github.com/YaleOpenLab/openx/rpc"
	utils "github.com/YaleOpenLab/openx/utils"
)

// server starts a local server which would inform us about the uptime of the teller and provide a data endpoint
type Data struct {
	// the data that is oging to be streamed
	// TODO: define what goes in here
	Timestamp string
	Info      string
}

func checkGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "404 page not found", http.StatusNotFound)
	}
}

func responseHandler(w http.ResponseWriter, r *http.Request, status int) {
	var response rpc.StatusResponse
	response.Code = status
	switch status {
	case rpc.StatusOK:
		response.Status = "OK"
	case rpc.StatusBadRequest:
		response.Status = "Bad Request error!"
	case rpc.StatusNotFound:
		response.Status = "404 Error Not Found!"
	case rpc.StatusInternalServerError:
		response.Status = "Internal Server Error"
	}
	rpc.MarshalSend(w, r, response)
}

// setupDefaultHandler sets up the default handler (ie returns 404 for invalid routes)
func setupDefaultHandler() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		responseHandler(w, r, rpc.StatusNotFound)
		return
	})
}

func pingHandler() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		var pr rpc.StatusResponse
		pr.Code = 200
		prJson, err := json.Marshal(pr)
		if err != nil {
			responseHandler(w, r, rpc.StatusInternalServerError)
			return
		}
		WriteToHandler(w, prJson)
	})
}

// TODO: read the data from the zigbee devices here
// so that we can verify the untrusted certificate.
// also clients who want this information can use this API directly without
// requiring a streaming service to inform them about changes. The client
// can call the teller and ask for data instantly and the API will respond.

// one problem here is to serve as less data as possible simply because
// we will be retruning mroe data thatn we receive and this becomes a trivial data
// amplification attack. IN this case, it is better to serve an ipfs hash instead
// of serving the whole blob of data
func dataHandler() {
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		var topsecret Data
		topsecret.Timestamp = utils.Timestamp()
		topsecret.Info = "this data is top secret and is for eyes only"
		hash, err := GetIpfsHash(topsecret.Info)
		if err != nil {
			log.Println(err)
		}
		log.Println("IPFS HASH IS: ", hash)
		// this is the data we need to pull in from the zigbee devices
		topsecretJson, err := json.Marshal(topsecret)
		if err != nil {
			responseHandler(w, r, rpc.StatusInternalServerError)
			return
		}
		WriteToHandler(w, topsecretJson)
	})
}

func SetupRoutes() {
	setupDefaultHandler()
	pingHandler()
	dataHandler()
}

// curl https://localhost/ping --insecure {"Code":200,"Status":""}
func StartServer() {
	SetupRoutes()
	err := http.ListenAndServeTLS(":"+consts.Tlsport, "ssl/server.crt", "ssl/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

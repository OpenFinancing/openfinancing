package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	consts "github.com/OpenFinancing/openfinancing/consts"
	rpc "github.com/OpenFinancing/openfinancing/rpc"
	utils "github.com/OpenFinancing/openfinancing/utils"
)

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

func checkPost(w http.ResponseWriter, r *http.Request) {
	log.Println("Checking POST")
	if r.Method != "POST" {
		http.Error(w, "404 page not found", http.StatusNotFound)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	// TODO, tweak this part
	if status == http.StatusNotFound {
		fmt.Fprint(w, "Invalid request params")
	}
}

func PingHandler() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		var pr rpc.PingResponse
		pr.Status = "Alive"
		prJson, err := json.Marshal(pr)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, prJson)
	})
}

func DataHandler() {
	// TODO: we need to read the data from the zigbee devices here
	// also clients which want this information can use the API directly without
	// the teller requiring a streaming service to inform them about changes. The client
	// can call the teller and ask for data at the instant and  the API should respond.
	// Takes less energy on the teller (which will be running on a low powered device)
	// and also saves a ton of complexity on our side. Also, the  cert gives us
	// ssl, so no mitm, which should alleviate problems arising from streaming.
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		var topsecret Data
		topsecret.Timestamp = utils.Timestamp()
		topsecret.Info = "this data is top secret and is for eyes only"
		topsecretJson, err := json.Marshal(topsecret)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, topsecretJson)
	})
}

func SetupRoutest() {
	PingHandler()
	DataHandler()
}

func StartServer() {
	SetupRoutest()
	err := http.ListenAndServeTLS(":"+consts.Tlsport, "ssl/server.crt", "ssl/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

package rpc

// the rpc package contains functions related to the server which will be interacting
// with the frontend. Not expanding on this too much since this will be changing quite often
// TODO: update RPC package in line with recent changes on master
// also evaluate on how easy it would be to rewrite this in nodeJS since the
// frontend is in react. Not many advantages per se and this works fine, so I guess
// we'll stay with this one for a while
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PingResponse struct {
	// pingresponse returns "alive" when called could be used by services
	// that scan for uptime
	Status string
}

type StatusResponse struct {
	Status int
}

// setupBasicHandlers sets up two hadnler functions that can be used to serve a default
// 404 response when we either error out or received input is incorrect.  This is not
// exactly ideal, because we don't expcet the RPC to be exposed and would like some more
// errors when we handle it on the frontend, but this makes for more a bit more
// secure Frontedn implementation which doesn't leak any information to the frontend
func setupBasicHandlers() {
	setupDefaultHandler()
	setupPingHandler()
}

func WriteToHandler(w http.ResponseWriter, jsonString []byte) {
	w.Header().Add("Access-Control-Allow-Origin", "localhost")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonString)
}

func checkOrigin(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "localhost" { // allow only our frontend UI to connect to our RPC instance
		http.Error(w, "404 page not found", http.StatusNotFound)
	}
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
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 Page Not Found")
	}
}

func setupDefaultHandler() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// default to 404 for every application not running on localhost
		errorHandler(w, r, http.StatusNotFound)
		return
	})
}

func setupPingHandler() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkGet(w, r)
		var pr PingResponse
		pr.Status = "Alive"
		prJson, err := json.Marshal(pr)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, prJson)
	})
}

// StartServer runs on the server side ie the server with the frontend.
// having to define specific endpoints for this because this
// is the system that would be used by the backend, so has to be built secure.
func StartServer(port string) {
	// we have a couple sub handlers for each main handler. these handlers
	// call the relevant internal endpoints and return a status / data.
	// we also have to process data from the pi itself, and that should have its own
	// functions somewhere else that can be accessed by the rpc.

	// also, this is assumed to run on localhost and hence has no authentication mehcanism.
	// in the case we want to expose the API, we must add some stuff that secures this.
	// right now, its just the CORS header, since we want to allwo all localhost processes
	// to access the API
	// a potential improvement will be to add something like macaroons
	// so that we can serve over an authenticated channel
	setupBasicHandlers()
	// setup basic handlers - / and /ping
	setupProjectRPCs()
	// setup project related RPCs
	setupInvestorRPCs()
	// setup investor related RPCs
	setupRecipientRPCs()
	// setup recipient related RPCs
	setupBondRPCs()
	setupCoopRPCs()
	portString := ":" + port // weird construction, but this should work
	log.Fatal(http.ListenAndServe(portString, nil))
}

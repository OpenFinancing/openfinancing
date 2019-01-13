package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"log"
	"os"
	"os/signal"
	"strings"

	consts "github.com/OpenFinancing/openfinancing/consts"
)

// package teller contains the remote client code that would be run on the client's
// side and communicate information with us and with atonomi and other partners.
// for now, we need a client that can start up and generate a pk/seed pair. This would
// be stored in the project struct and if anyone wants to see the status of this
// node, they can check the blockchain for the node's updates (node should update
// the blockchain at frequient intervals with power generation data in the memo
// field of the tx.
// Polling interval would be 5 minutes, 1440/5 = 288 updates a day
// MW Is this the amount of calls from the rasberry pi per day? This might be too much.
// Not sure what the colors here mean.
// SHould this be what is listening to a protected MQTT channel?

var PublicKey string
var Seed string

func main() {
	CreateHomeDir()
	CreateFile()
	Authenticate()
	promptColor := color.New(color.FgHiYellow).SprintFunc()
	whiteColor := color.New(color.FgHiWhite).SprintFunc()
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      promptColor("teller") + whiteColor("# "),
		HistoryFile: consts.TellerHomeDir + "/history.txt",
		// AutoComplete: lc.NewAutoCompleter(),
	})

	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()
	// main shell loop
	var t Client
	t.Info = "Console CLI Client for testing"
	t.Location = "Location unknown. Scanning.."
	t.UniqueId = "password" // the thing we need for unlocking the auth.txt file
	go StartServer()
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("\nSigint received in quit function. not quitting!")
		close(cleanupDone)
	}()
	for {
		// setup reader with max 4K input chars
		msg, err := rl.Readline()
		if err != nil {
			var err error
			err = EndHandler(t) // error, user wants to quit
			for err != nil {
				log.Println(err)
				err = EndHandler(t)
				<-cleanupDone // to prevent user from quitting when sigin arrives
			}
			break
		}
		msg = strings.TrimSpace(msg)
		if len(msg) == 0 {
			continue
		}
		rl.SaveHistory(msg)

		cmdslice := strings.Fields(msg)
		ColorOutput("entered command: "+msg, YellowColor)

		err = ParseInput(cmdslice)
		if err != nil {
			fmt.Println(err)
		}
	}
}

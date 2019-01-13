package main

import (
	"log"

	"github.com/Varunram/gobee" // importing fork because package is not maintained anymore and might be altered
	"github.com/Varunram/gobee/api/rx"
)

// this implements the xbee receiving stuff from the iot devices

type Receiver struct {
	Receiver gobee.XBeeReceiver
}

func (*Receiver) Receive(frame rx.Frame) error {
	switch frame.(type) {
	case *rx.ZB:
		log.Println("Received frame!")
		// do other stuff here, can;t really know what without testing with a zigbee module
	}
	return nil
}

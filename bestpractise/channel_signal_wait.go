package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	// monitor keyboardinterrupt.
	var sig = make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	for {
		fmt.Println("........")
		time.Sleep(200)
		// chk(binary.Write(f, binary.BigEndian, in))
		// nSamples += len(in)
		select {
		case <-sig:
			return
		default:
		}
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/schoentoon/piglow"
)

var colorStr = flag.String("color", "red", "color to flash")
var interval = flag.Duration("interval", time.Second*3, "interval length")

func main() {
	flag.Parse()
	if ok := piglow.HasPiGlow(); !ok {
		log.Fatalf("piglow not available")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Printf("got sig kill - shutting down")
		err := piglow.ShutDown()
		if err != nil {
			log.Printf("failed to shutdown piglow: %v", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	var color byte
	switch strings.ToLower(*colorStr) {
	case "red":
		color = piglow.Red
	case "orange":
		color = piglow.Orange
	case "yellow":
		color = piglow.Yellow
	case "green":
		color = piglow.Green
	case "blue":
		color = piglow.Blue
	case "white":
		color = piglow.White
	}

	var on bool
	var intensity byte
	for {
		intensity = 0
		if !on {
			intensity = 255
		}
		on = !on
		piglow.Ring(color, intensity)
		time.Sleep(*interval)
	}
}

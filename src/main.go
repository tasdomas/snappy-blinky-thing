package main

import (
	"log"
	"time"

	"github.com/schoentoon/piglow"
)

func main() {

	if ok := piglow.HasPiGlow(); !ok {
		log.Fatalf("piglow not available")
	}

	var on bool
	for {
		intensity := byte(0)
		if on {
			on = false
		} else {
			intensity = 255
			on = true
		}
		err := piglow.Ring(piglow.Red, intensity)
		if err != nil {
			log.Fatalf("failed to adjust led status: %v", err)
		}
		time.Sleep(time.Second)
	}
}

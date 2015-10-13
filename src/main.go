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

	legs := make([]bool, 3)
	var i byte
	var err error
	var intensity byte
	for {
		if legs[i] {
			intensity = 0
			legs[i] = false
		} else {
			intensity = 255
			legs[i] = true
		}
		err = piglow.Leg(i, intensity)
		if err != nil {
			log.Fatalf("failed to adjust led status: %v", err)
		}
		i = (i + 1) % 3
		time.Sleep(time.Second)
	}
}

package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/schoentoon/piglow"
)

const (
	BLINK_INTERVAL      = time.Second * 5
	INTENSITY_MAX  byte = 4

	StatusURL     = "http://juju.fail/status.json"
	Refresh       = time.Minute
	WatchedBranch = "master"

	APP_PATH_KEY = "SNAPP_APP_PATH"
)

// do_blink listens on the provided channel for specified color and blinks it at an
// interval.
func do_blink(color <-chan byte) {
	var c byte
	var intensity byte

	tick := time.NewTicker(BLINK_INTERVAL)

LOOP:
	for {
		select {
		case c_new, ok := <-color:
			if !ok {
				break LOOP
			}
			piglow.Ring(c, 0)
			c = c_new
		case <-tick.C:
			if intensity == 0 {
				intensity = INTENSITY_MAX
			} else {
				intensity = 0
			}
			piglow.Ring(c, intensity)
		}
	}
	err := piglow.ShutDown()
	if err != nil {
		log.Printf("failed to shutdown piglow: %v", err)
	}
}

// JujuFail decodes the information as it is presented in
// http://juju.fail/status.json.
type JujuFail struct {
	Status  map[string][]interface{}
	Updated string
}

func get_juju_status() (JujuFail, error) {
	res, err := http.Get(StatusURL)
	if err != nil {
		log.Printf("failed to query URL %q: %v", StatusURL, err)
		return JujuFail{}, err
	}
	dec := json.NewDecoder(res.Body)

	var current JujuFail
	err = dec.Decode(&current)
	res.Body.Close()
	if err != nil && err != io.EOF {
		log.Printf("could not decode juju status: %v", err)
		return JujuFail{}, err
	}
	return current, nil
}

func main() {
	var knownStatus *JujuFail

	if ok := piglow.HasPiGlow(); !ok {
		log.Fatalf("piglow not available")
	}

	colors := make(chan byte)
	go do_blink(colors)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	interval := time.NewTicker(Refresh)
MAINLOOP:
	for {
		select {
		case <-interval.C:
			colors <- piglow.Orange

			current, err := get_juju_status()
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}
			// Ignore unchanged status.
			if knownStatus != nil {
				if len(current.Status[WatchedBranch]) == len(knownStatus.Status[WatchedBranch]) {
					continue
				}
			}

			color := piglow.Green
			if len(current.Status[WatchedBranch]) != 0 {
				color = piglow.Red
			}
			colors <- color
			knownStatus = &current
		case <-sigChan:
			break MAINLOOP
		}
	}

	log.Printf("got sig kill - shutting down")
	err := piglow.ShutDown()
	if err != nil {
		log.Printf("failed to shutdown piglow: %v", err)
		os.Exit(1)
	}
	os.Exit(0)

}

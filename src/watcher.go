package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"
)

const (
	StatusURL     = "http://juju.fail/status.json"
	Refresh       = time.Minute
	WatchedBranch = "master"

	APP_PATH_KEY = "SNAPP_APP_PATH"
)

// JujuFail decodes the information as it is presented in
// http://juju.fail/status.json.
type JujuFail struct {
	Status  map[string][]interface{}
	Updated string
}

func changeBlinking(current *os.Process, color string, interval int) (*os.Process, error) {
	if current != nil {
		err := current.Kill()
		if err != nil {
			log.Printf("failed to kill current blinker: %v", err)
			return nil, err
		}
		_, err = current.Wait()
		if err != nil {
			log.Printf("failed waiting for  current blinker to die: %v", err)
		}
	}
	base := os.Getenv(APP_PATH_KEY)
	blinker := path.Join(base, "bin", "blinker")
	cmd := exec.Command(blinker, "-color", color, "-interval", fmt.Sprintf("%ds", interval))
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	p, err := os.FindProcess(cmd.Process.Pid)
	if p == nil || err != nil {
		return nil, err
	}
	return cmd.Process, nil
}

func main() {
	var blinker *os.Process
	var err error
	var knownStatus *JujuFail

	for {
		time.Sleep(Refresh)
		blinker, err = changeBlinking(blinker, "YELLOW", 1)
		if blinker == nil {
			log.Printf("failed to start blinker: %v", err)
		}
		res, err := http.Get(StatusURL)
		if err != nil {
			log.Printf("failed to query URL %q: %v", StatusURL, err)
			continue
		}
		dec := json.NewDecoder(res.Body)

		var current JujuFail
		err = dec.Decode(&current)
		res.Body.Close()
		if err != nil && err != io.EOF {
			log.Printf("could not decode juju status: %v", err)
			continue
		}

		// Ignore unchanged status.
		if knownStatus != nil {
			if len(current.Status[WatchedBranch]) == len(knownStatus.Status[WatchedBranch]) {
				continue
			}
		}

		color := "GREEN"
		if len(current.Status[WatchedBranch]) != 0 {
			color = "RED"
		}
		blinker, err = changeBlinking(blinker, color, 1)
		if blinker == nil {
			log.Printf("failed to start blinker: %v", err)
		} else {
			knownStatus = &current
		}
	}
}

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

const (
	ntpServerURL  = "0.beevik-ntp.pool.ntp.org"
	timePrecision = 1 * time.Second
)

func main() {
	currentTime := time.Now()
	exactTime, err := ntp.Time(ntpServerURL)

	if err != nil {
		log.Fatalf("Error querying time from a remote NTP server \"%s\": %v\n", ntpServerURL, err)
	}

	fmt.Printf("current time: %v\n", currentTime.Round(timePrecision))
	fmt.Printf("exact time: %v\n", exactTime.Round(timePrecision))
}

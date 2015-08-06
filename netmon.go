package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var (
	netmon = newNetMon()
)

func init() {
	go netmon.Run()
}

type netMon struct {
	Frequency     time.Duration
	RateRX        uint64
	RateTX        uint64
	Interfaces    []string
	CounterActive bool

	lastRX uint64
	lastTX uint64
}

func newNetMon() *netMon {
	return &netMon{
		Frequency:     time.Second,
		Interfaces:    []string{},
		CounterActive: true,
	}
}

func (n *netMon) Run() {
	for n.CounterActive {
		n.update()
		<-time.After(n.Frequency)
	}
}

func (n *netMon) update() {
	netContent, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		n.CounterActive = false
		return
	}

	var (
		currentRX uint64
		currentTX uint64
	)
	for _, line := range strings.Split(string(netContent), "\n") {
		if !strings.Contains(line, ":") {
			continue
		}

		fields := n.removeEmpty(strings.Split(line, " "))
		for _, intf := range n.Interfaces {
			if intf == strings.Trim(fields[0], ":") {
				rx, err := strconv.ParseUint(fields[1], 10, 64)
				if err != nil {
					fmt.Printf("Unable to parse RX for IF %s: %s\n", intf, err)
					return
				}
				tx, err := strconv.ParseUint(fields[9], 10, 64)
				if err != nil {
					fmt.Printf("Unable to parse TX for IF %s: %s\n", intf, err)
					return
				}

				currentRX += rx
				currentTX += tx
			}
		}
	}

	n.RateRX, n.lastRX = currentRX-n.lastRX, currentRX
	n.RateTX, n.lastTX = currentTX-n.lastTX, currentTX
}

func (n *netMon) removeEmpty(in []string) (out []string) {
	for i := range in {
		if len(in[i]) > 0 {
			out = append(out, in[i])
		}
	}

	return
}

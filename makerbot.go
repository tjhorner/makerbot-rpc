// Package makerbot is a Go client library for MakerBot printers.
//
// Full README: https://github.com/tjhorner/makerbot-rpc/blob/master/README.md
package makerbot

import (
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/mdns"
)

// These constants are used to communicate with the printer
// and are apparently hard-coded

const (
	makerbotClientID     = "MakerWare"
	makerbotClientSecret = "secret"
)

// NewClient creates a new client.
func NewClient() Client {
	return Client{Connected: false, Timeout: 5 * time.Second}
}

const mdnsService = "_makerbot-jsonrpc._tcp"

// DiscoverPrinters will discover printers that are on your current
// LAN network and return them once `timeout` is up.
//
// Note that all fields are not returned in the printer's mDNS TXT
// reply. Useful fields like MachineName and IP/Port are in there,
// though, so that should be enough to initiate a connection with
// the printer.
func DiscoverPrinters(timeout time.Duration) (*[]Printer, error) {
	var printers []Printer

	ch := make(chan *mdns.ServiceEntry)
	go func() {
		for entry := range ch {
			if !strings.Contains(entry.Name, "_makerbot-jsonrpc") {
				continue
			}

			fields := *parseInfoFields(&entry.InfoFields)

			vid, _ := strconv.Atoi(fields["vid"])
			pid, _ := strconv.Atoi(fields["pid"])

			printer := Printer{
				MachineName:        fields["machine_name"],
				MachineType:        fields["machine_type"],
				APIVersion:         fields["api_version"],
				Serial:             fields["iserial"],
				MotorDriverVersion: fields["motor_driver_version"],
				Vid:                vid,
				Pid:                pid,
				SSLPort:            fields["ssl_port"],
				BotType:            fields["bot_type"],
				IP:                 entry.AddrV4.String(),
				Port:               string(entry.Port),
			}

			printers = append(printers, printer)
		}
	}()

	params := mdns.DefaultParams(mdnsService)
	params.Timeout = timeout
	params.Entries = ch

	err := mdns.Query(params)
	if err != nil {
		return nil, err
	}

	close(ch)

	return &printers, nil
}

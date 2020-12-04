package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type app struct {
	gopi.Unit
	gopi.Metrics
	gopi.MetricWriter
	gopi.Logger
	gopi.Command
	*http.Client

	timeout     *time.Duration
	delta       *time.Duration
	measurement *string
	ip          net.IP
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	uriApify  = "https://api.ipify.org"
	uriGoogle = "https://domains.google.com/nic/update"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *app) Define(cfg gopi.Config) error {
	// Register commands
	cfg.Command("ip", "Display External IP Address", this.Once)
	cfg.Command("daemon", "Start daemon", this.Daemon)

	// Register flags
	this.delta = cfg.FlagDuration("delta", 5*time.Minute, "Period to re-register when changed")
	this.timeout = cfg.FlagDuration("timeout", 15*time.Second, "Request timeout")
	this.measurement = cfg.FlagString("measurement", "dnsregister", "Measurement name")

	// Return success
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	if this.Command = cfg.GetCommand(cfg.Args()); this.Command == nil {
		return gopi.ErrHelp
	}

	// Set HTTP client
	this.Client = &http.Client{
		Timeout: *this.timeout,
	}

	// Create a measurement
	if *this.measurement != "" {
		host, err := os.Hostname()
		if err != nil {
			return err
		}
		tag := this.Metrics.Field("host", host)
		if _, err := this.Metrics.NewMeasurement(*this.measurement, "duration float64, subdomain,ip,status string", tag); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RUN

func (this *app) Run(ctx context.Context) error {
	return this.Command.Run(ctx)
}

func (this *app) Once(ctx context.Context) error {
	if ip, err := this.GetExternalAddress(); err != nil {
		return err
	} else {
		fmt.Println(ip)
	}

	// Return success
	return nil
}

func (this *app) Daemon(ctx context.Context) error {
	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	subdomain, user, passwd := GetCredentials()

	fmt.Println("Press CTRL+C to end")
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			now := time.Now()
			this.Debug("Discovery")
			if ip, err := this.GetExternalAddress(); err != nil {
				this.Print(err)
			} else if ip.Equal(this.ip) {
				this.Debug("...no change: ", this.ip)
				if err := this.Emit(time.Since(now).Seconds(), "", ip.String(), "nochg"); err != nil {
					this.Print(err)
				}
			} else if status, err := this.RegisterExternalAddress(ip, subdomain, user, passwd); err != nil {
				this.Print("Error", err)
			} else {
				this.Debug("...registration: ", ip, " ", status)

				// Update stored IP address when successful
				if status == "good" || status == "nochg" {
					this.ip = ip
				} else {
					this.ip = nil
				}

				// Metrics
				if err := this.Emit(time.Since(now).Seconds(), subdomain, ip.String(), status); err != nil {
					this.Print(err)
				}
			}

			// Reset timer to the next delta point
			timer.Reset(*this.delta)
		}
	}
}

func GetCredentials() (string, string, string) {
	subdomain, _ := os.LookupEnv("GOOGLE_DNS_SUBDOMAIN")
	user, _ := os.LookupEnv("GOOGLE_DNS_USER")
	passwd, _ := os.LookupEnv("GOOGLE_DNS_PASSWORD")
	return subdomain, user, passwd
}

func (this *app) Emit(latency float64, subdomain, ip, status string) error {
	if *this.measurement != "" {
		return this.Metrics.Emit(*this.measurement, latency, subdomain, ip, status)
	} else {
		return nil
	}
}

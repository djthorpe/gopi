package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/table"
)

////////////////////////////////////////////////////////////////////////////////

type hosts struct {
	gopi.ServiceRecord
}

type addrs struct {
	gopi.ServiceRecord
}

type txt struct {
	gopi.ServiceRecord
}

func (h hosts) Format() (string, table.Alignment, table.Color) {
	str := net.JoinHostPort(h.Host(), fmt.Sprint(h.Port()))
	return str, table.Auto, table.None
}

func (a addrs) Format() (string, table.Alignment, table.Color) {
	str := ""
	for i, a := range a.Addrs() {
		if i > 0 {
			str += " "
		}
		str += a.String()
	}
	return str, table.Auto, table.None
}

func (t txt) Format() (string, table.Alignment, table.Color) {
	str := ""
	for i, t := range t.Txt() {
		if i > 0 {
			str += " "
		}
		str += t
	}
	return str, table.Auto, table.None
}

////////////////////////////////////////////////////////////////////////////////

func (this *app) RunDiscovery(ctx context.Context) error {
	args := this.Command.Args()
	ctx, cancel := context.WithTimeout(ctx, *this.timeout)
	defer cancel()

	if this.ServiceDiscovery == nil {
		return gopi.ErrInternalAppError.WithPrefix("ServiceDiscovery")
	}

	if len(args) == 0 {
		return this.RunDiscoveryEnumerate(ctx)
	} else if len(args) == 1 {
		return this.RunDiscoveryLookup(ctx, args[0])
	} else {
		return gopi.ErrHelp
	}
}

func (this *app) RunDiscoveryEnumerate(ctx context.Context) error {

	// Enumerate services
	services, err := this.ServiceDiscovery.EnumerateServices(ctx)
	if err != nil {
		return err
	}

	// Display platform information
	table := table.New()
	table.SetHeader(header{"Service"})
	for _, service := range services {
		table.Append(service)
	}
	table.Render(os.Stdout)

	// Return success
	return nil
}

func (this *app) RunDiscoveryLookup(ctx context.Context, name string) error {
	// Enumerate services
	records, err := this.ServiceDiscovery.Lookup(ctx, name)
	if err != nil {
		return err
	}

	// Display service information
	table := table.New()
	table.SetHeader(header{"Service"}, header{"Name"}, header{"Host"}, header{"Addr"}, header{"Txt"})
	for _, record := range records {
		table.Append(record.Service(), record.Name(), hosts{record}, addrs{record}, txt{record})
	}
	table.Render(os.Stdout)

	// Return success
	return nil
}

func (this *app) RunDiscoveryServe(ctx context.Context) error {
	args := this.Command.Args()
	if this.ServiceDiscovery == nil {
		return gopi.ErrInternalAppError.WithPrefix("ServiceDiscovery")
	}

	record, err := this.GetServiceRecord(args)
	if err != nil {
		return err
	}

	// Display service information
	table := table.New()
	table.SetHeader(header{"Service"}, header{"Name"}, header{"Host"}, header{"Addr"}, header{"Txt"})
	table.Append(record.Service(), record.Name(), hosts{record}, addrs{record}, txt{record})
	table.Render(os.Stdout)

	fmt.Println("Serving, press CTRL+C to exit")

	// Serve until CTRL+C
	return this.ServiceDiscovery.Serve(ctx, []gopi.ServiceRecord{record})
}

func (this *app) GetServiceRecord(args []string) (gopi.ServiceRecord, error) {
	txt := []string{}
	service := "gopi"
	name := ""
	if hostname, err := os.Hostname(); err != nil {
		return nil, err
	} else {
		name = hostname
	}
	if *this.name != "" {
		name = strings.TrimSpace(*this.name)
	}
	if len(args) > 0 {
		service = args[0]
		txt = args[1:]
	}
	if strings.HasPrefix(service, "_") == false {
		service = "_" + service
	}
	if strings.HasSuffix(service, "._tcp") == false && strings.HasSuffix(service, "._udp") == false {
		service = service + "._tcp"
	}
	if record, err := this.ServiceDiscovery.NewServiceRecord(service, name, uint16(*this.port), txt, gopi.SERVICE_FLAG_IP4); err != nil {
		return nil, err
	} else {
		return record, nil
	}
}

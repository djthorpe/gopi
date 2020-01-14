
# Infrared Remote Control

The Linux Infrared Remote Control interface \(LIRC\) allows your hardware to send and receive IR codes and is compatible with many remote controls which use the IR protocol. More information is available from [kernel.org](https://www.kernel.org/doc/html/latest/media/uapi/rc/lirc-dev.html).

The LIRC Unit allows you to send and receive information from LIRC where you have installed an IR sender and receiver device. In order to use the Unit, here are the parameters:

{% hint style="info" %}

| Parameter | Value |
| :--- | :--- |
| Name | `gopi/lirc` |
| Interface | `gopi.LIRC` |
| Type | `gopi.UNIT_LIRC` |
| Import | `github.com/djthorpe/gopi/v2/unit/lirc` |
| Events | `gopi.LIRCEvent` |
| Compatibility | Linux |

{% endhint %}

The unit adheres to the following interface:

```go
type gopi.LIRC interface {
    // Modes
    RcvMode() gopi.LIRCMode
	SendMode() gopi.LIRCMode
	SetRcvMode(gopi.LIRCMode) error
	SetSendMode(gopi.LIRCMode) error

	// Receive parameters
	RcvDutyCycle() uint32
	GetRcvResolution() (uint32, error)
	SetRcvTimeout(micros uint32) error
	SetRcvTimeoutReports(enable bool) error
	SetRcvCarrierHz(value uint32) error
	SetRcvCarrierRangeHz(min uint32, max uint32) error

	// Send parameters
	SendDutyCycle() uint32
	SetSendDutyCycle(value uint32) error
	SetSendCarrierHz(value uint32) error

	// Send Pulse mode, values are in milliseconds
	PulseSend(values []uint32) error

	// Implements gopi.Unit
	gopi.Unit
}
```

Here is a short example of how to receive IR codes from a remote control:

```go
package main
// ...
var (
	Events = []gopi.EventHandler{ gopi.EventHandler{Name: "gopi.LIRCEvent", Handler: LIRCHandler} }
)

func LIRCHandler(_ context.Context, _ gopi.App, evt_ gopi.Event) {
	evt := evt_.(gopi.LIRCEvent)
	fmt.Printf("%-10s %-10s %v\n",evt.Mode(),evt.Type(),evt.Value())
}

func Main(app gopi.App, args []string) error {
    app.WaitForSignal(context.Background(), os.Interrupt)
	return nil
}

func main() {
	if app, err := app.NewCommandLineTool(Main, Events, "lirc"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		os.Exit(app.Run())
	}
}
```

Typical output might look like this:
```bash

MODE       TYPE       VALUE
---------- ---------- --------------------
mode2      space      20812
mode2      pulse      2636
mode2      space      3332
mode2      pulse      335
...
mode2      space      1576
mode2      timeout    22605
```

Any LIRC signal received when the receive uses `LIRC_MODE_MODE2` is a series of pulse and space values, measured in microseconds. If your LIRC driver supports timeouts, you may also receive timeout messages.
Ultimately an IR signal is a digital waveform of varying "on" and "off" signals. In `LIRC_MODE_MODE2`
your code will need to decode these into control (or "scan" codes). There are various protocols to do
this, and these are documented [here](https://www.kernel.org/doc/html/latest/media/uapi/rc/lirc-dev-intro.html#lirc-modes).

## Hardware set-up

You'll need an IR sender and/or receiver plugged into your GPIO port. One example is the Energenie PiMote [link](https://energenie4u.co.uk/catalogue/product/ENER314-IR) but you can also construct your own from a PCB, one example is [here](github.com/djthorpe/remotes).

For a Raspberry Pi, you should enable the LIRC drivers in your kernel by adding the following line to your `/boot/config.txt` file modifying the GPIO pins in order to load the LIRC (Linux Infrared Control) driver, and then reboot your Raspberry Pi:

```bash
dtoverlay=gpio-ir,gpio_pin=18
dtoverlay=gpio-ir-tx,gpio_pin=17
```

On older kernels there is a single kernel driver:

```bash
dtoverlay=lirc-rpi,gpio_in_pin=18,gpio_out_pin=17
```

You should then be able to see the device drivers `/dev/lirc0` and `/dev/lirc1`. If not check output of the lsmod command. You will also need to give your user read and write access to the device drivers:

```
and ensure your user has the correct permissions to access the device using the following command:

```bash
bash% sudo usermod -a -G video ${USER}
```

## Receiving & decoding IR signals

_Section to be written_

## Sending IR signals

_Section to be written_

## Example application

_Section to be written_










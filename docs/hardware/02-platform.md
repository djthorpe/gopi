---
description: Obtaining information about your hardware platform.
---

# Platform Information

The **Platform Unit** returns some information about the platform your tool is running on.

{% hint style="info" %}
| Parameter | Value |
| :--- | :--- |
| Name | `gopi/platform` |
| Interface | `gopi.Platform` |
| Type | `gopi.UNIT_PLATFORM` |
| Import | `github.com/djthorpe/gopi/v2/unit/platform` |
| Compatibility | Linux, Darwin, Raspberry Pi |
{% endhint %}

Here is the interface which the platform unit adheres to:

```go
type gopi.Platform interface {
    gopi.Unit

    Product() string     // Product returns product name
    Type() gopi.PlatformType     // Type returns flags identifying platform type
    SerialNumber() string     // SerialNumber returns unique serial number for host
    Uptime() time.Duration     // Uptime returns uptime for host
    LoadAverages() (float64, float64, float64)     // LoadAverages returns 1, 5 and 15 minute load averages
    NumberOfDisplays() uint     // NumberOfDisplays returns the number of possible displays for this host
}
```

Here's an example of accessing the platform information in your `Main` function:

```go
func Main(app gopi.App, args []string) error {
    platform := app.Platform()
    fmt.Println(platform.Type(),platform.Product(),platform.SerialNumber())
    // ...
}
```

There are some platform differences with the information returned:

* On Linux, the generic name `linux` is returned for product;
* On Linux, a Mac Address is returned for the serial number;
* On Darwin, the product is a product code \(ie, "MacPro1,2"\) rather than name.
* On Darwin and Linux, the number of displays is returned as zero as these platform displays are not yet supported.


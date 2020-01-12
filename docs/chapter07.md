
---
description: 'Accessing information about your display.'
---

# Displays

The Display Unit returns some information about your display. When importing this unit into your tool, the command line flag `-display` can be used to choose the display. An error will be returned when trying to use this unit on Linux or Darwin when the tool is run.

{% hint style="info" %}

| Parameter | Value |
| :--- | :--- |
| Name | `gopi/display` |
| Interface | `gopi.Display` |
| Type | `gopi.UNIT_DISPLAY` |
| Import | `github.com/djthorpe/gopi/v2/unit/display` |
| Compatibility | Raspberry Pi |

{% endhint %}

Here's an example of returning information about the display:

```go
func Main(app gopi.App, args []string) error {
    display := app.Display()
    fmt.Println(display.Name(),display.Size())
    // ...
}
```

Here is the interface for a display:

```go
type gopi.Display interface {
    gopi.Unit

    DisplayId() uint     // Return display number
    Name() string // Return name of the display
    Size() (uint32, uint32) // Return display size for nominated display number
    PixelsPerInch() uint32 // Return the PPI (pixels-per-inch) for the display
}
```


# rpi

## Creating a Raspberry Pi object

In order to control the Raspberry Pi, you'll need to create a device object
as follows:

```go

import "gopi/rpi"

func main() {
  device, err := rpi.New()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error: ", err)
    os.Exit(-1)
  }
  defer device.Close()
}
```
You can then query the Raspberry Pi for various information, for example:

```go
  warranty, err := device.WarrantyBit()
  product, err := device.Product() /* rpi.Product */
  processor, err := device.Processor() /* Processor */
  product_name, err := device.ProductName()
  processor_name, err := device.ProcessorName()
  peripheral_base, err := device.PeripheralBase()
```

Here are the current possible values for Product and Processor:

```go
	RPI_MODEL_UNKNOWN
	RPI_MODEL_A
	RPI_MODEL_B
	RPI_MODEL_A_PLUS
	RPI_MODEL_B_PLUS
	RPI_MODEL_B_PI_2
	RPI_MODEL_B_PI_3
	RPI_MODEL_ALPHA
	RPI_MODEL_COMPUTE_MODULE
	RPI_MODEL_ZERO
```

```go
	RPI_PROCESSOR_UNKNOWN
	RPI_PROCESSOR_BCM2835
	RPI_PROCESSOR_BCM2836
	RPI_PROCESSOR_BCM2837
```





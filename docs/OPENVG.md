
# Using OpenVG

The Raspberry Pi can render 2D vector graphics onto EGL surfaces with
GPU acceleration. As Wikipedia says,

| OpenVG is an API designed for hardware-accelerated 2D vector graphics. 
| Its primary platforms are mobile phones, gaming & media consoles and 
| consumer electronic devices.

## Abstract Interface

There are currently three main concepts which allow you to draw with OpenVG:

  * The **Driver** provides information about the surface you are drawing on,
    and allows you to create path and paint objects.
  * The **Path** object stores a set of points or shapes to draw on the surface,
    and allows you to stroke and fill the path.  
  * The **Paint** object stores the drawing state (color, line width, etc)
    with which to stroke and fill paths.

There are other enumerations, interfaces and structs which support these concepts:	

| **Import** | `github.com/djthorpe/gopi/khronos` |
| -- | -- | -- |
| **Interface** | `khronos.VGDriver` | gopi.Driver, the OpenVG driver |
| **Interface** | `khronos.VGPath` | The path |
| **Interface** | `khronos.VGPaint` | The paintbrush |
| **Struct** | `khronos.VGColor` | Painting colour |
| **Struct** | `khronos.VGPoint` | A point or size |

## Concrete Implementation

The concrete implementation for the Raspberry Pi requires the following
configuration:

| **Import** | `github.com/djthorpe/gopi/device/rpi` |
| -- | -- | -- |
| **Struct** | `rpi.OpenVG` | Concrete Raspberry Pi OpenVG driver |

The one argument to the configuration object is the "EGL" instance, which
should also be a concrete Raspberry Pi object. For example,

```go
  openvg, err := gopi.Open(rpi.OpenVG{ EGL: egl_object },logger)
  if err != nil { /* handle error */ }
  defer openvg.Close()

  // Clear to white
  openvg.(khronos.VGDriver).Begin()
  openvg.(khronos.VGDriver).Clear(khronos.VGColorWhite)
  openvg.(khronos.VGDriver).Flush()
```

The return value of `gopi.Open` is a `gopi.Driver` so in order to call driver
methods you need to cast to a `khronos.VGDriver`.

## The Coordinate system

Like the EGL interface, the co-ordinate system for OpenVG has the origin pixel
`khronos.VGPoint{ 0, 0 }` at the top left of the drawable surface.

TODO: Calculating points, alignment, etc.

## Drawing on surfaces

## Creating and Drawing Paths

## Paintbrushes


# Links

  * (Wikipedia Entry](https://en.wikipedia.org/wiki/OpenVG)
  * [A Technical Introduction to OpenVG](https://www.khronos.org/assets/uploads/developers/library/siggraph2006/OpenKODE_Course/05_OpenVG-Overview.pdf)
  

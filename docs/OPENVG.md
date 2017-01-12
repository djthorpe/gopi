
# Using OpenVG

The Raspberry Pi can render 2D vector graphics onto EGL surfaces with
GPU acceleration. As Wikipedia says,

> OpenVG is an API designed for hardware-accelerated 2D vector graphics. 
> Its primary platforms are mobile phones, gaming & media consoles and 
> consumer electronic devices.

The following sections explain the concepts behind the OpenVG
implemenation. 

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
  egl := /* EGL object */
  surface := /* Surface on which to draw */	
  openvg, err := gopi.Open(rpi.OpenVG{ EGL: egl },logger)
  if err != nil { /* handle error */ }
  defer openvg.Close()

  // Clear surface to white
  openvg.(khronos.VGDriver).Begin(surface)
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

In order to draw on an EGL surface, you will need to create a surface which is
"bound" to the OpenVG API. For example,

```go
  surface, err := app.EGL.CreateBackground("OpenVG",1.0)
  if err != nil {
    return err
  }
  defer app.EGL.DestroySurface(surface)
```

Painting is an **atomic** operation. You should call the driver `Begin` method
to start drawing and the `Flush` method to end the drawing. Once `Flush` is 
called, the drawing is made visible. A syncronization lock will ensure you
cannot call `Begin` more than once without flushing. If you don't want to
lock then you should use `BeginNoWait` instead, which will return an error
immediately if the drawing surface is already locked.

For example, here's a wrapper function which could ensure the `Begin`/`Flush`
methods are always used together:

```go
  func draw_surface(openvg khronos.VGDriver,surface khronos.EGLSurface,callback func() error) error {
    if err := openvg.Begin(surface); err != nil {
  	  return error
    }
    defer openvg.Flush()
	return callback()
  }
```

## Creating and Drawing Paths

## Paintbrushes


# Links

  * (Wikipedia Entry](https://en.wikipedia.org/wiki/OpenVG)
  * [A Technical Introduction to OpenVG](https://www.khronos.org/assets/uploads/developers/library/siggraph2006/OpenKODE_Course/05_OpenVG-Overview.pdf)
  

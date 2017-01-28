
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
```

The return value of `gopi.Open` is a `gopi.Driver` so in order to call driver
methods you need to cast to a `khronos.VGDriver`.

## The Coordinate system

Like the EGL interface, the co-ordinate system for OpenVG has the origin pixel
`khronos.VGPoint{ 0, 0 }` at the top left of the drawable surface. Unlike the
EGL interface, points are defined as the `float32` type rather than `int`,
so that points can be transformed (rotated, etc) by fractional amounts.

TODO: Calculating points, alignment, etc.

The following methods can be used for "transforming" the co-ordinate system
when drawing:

| **Interface** | `khronos.VGDriver` |
| -- | -- | -- |
| **Method** | `Translate(offset VGPoint) error` | Translate co-ordinate sysyetm to offset drawing |
| **Method** | `Scale(x,y float32) error` | Scale co-ordinate system |
| **Method** | `Shear(x,y float32) error` | Shear co-ordinate system |
| **Method** | `Rotate(r float32) error` | Rotate co-ordinate system |
| **Method** | `LoadIdentity() error` | Reset co-ordinate system |

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

Painting is an **atomic** operation. As such, you draw by calling the `Do` method
with an argument to your drawing callback. Only one `Do` method can be called
at any one time:

```go
  surface := /* Surface on which to draw */	
  openvg.(khronos.VGDriver).Do(surface,func () error {
	/* Draw on surface */
	openvg.(khronos.VGDriver).Clear(surface,khronos.VGColorWhite)
	
	/* Return success */
	return nil
  });
  if err != nil {
    return err
  }
```

Before the drawing begins, the point transformation identity matrix is loaded,
so that there is a one-to-one correlation between points and pixels. You can then
rotate or transform the co-ordinate system before you start drawing. After the
drawing is completed, the surface is flushed.

The `khronos.VGDriver` interface implements the following methods:

| **Interface** | `khronos.VGDriver` |
| -- | -- | -- |
| **Method** | `Do(surface EGLSurface, callback func() error) error` | Draw on surface |
| **Method** | `Clear(surface EGLSurface, color VGColor) error` | Clear surface to color |

The `Clear` method will clear a surface to the specified color, without taking note of
any co-ordinate transformations.

## Paintbrushes

A paintbrush is used for drawing outlines (known as the "stroke") and filling
shapes (known as the "path"). In order to create a paintbrush, the `khronos.VGDriver`
interface implements the following methods:

| **Interface** | `khronos.VGDriver` |
| -- | -- | -- |
| **Method** | `CreatePaint(color VGColor) (VGPaint, error)` | Create a paintbrush for stroking or filling |
| **Method** | `DestroyPaint(VGPaint) error` | Destroy a created paintbrush |

When creating a paintbrush, the colour of the paintbrush is provided and
the following defaults are set:

  * `VGFillRule` is set to `VG_STYLE_FILL_EVENODD`
  * `VGStrokeWidth` is set to 1.0
  * `VGStrokeCapStyle` is set to `VG_STYLE_CAP_BUTT`
  * `VGStrokeJoinStyle` is set to `VG_STYLE_JOIN_MITER`
  
You can affect the paintbrush attributes using the following methods:

| **Interface** | `khronos.VGPaint` |
| -- | -- | -- |
| **Method** | `SetColor(color VGColor) error` | Set paintbrush colour |
| **Method** | `SetFillRule(style VGFillRule) error` | Set fill rule (for filling) |
| **Method** | `SetStrokeWidth(width float32) error` | Set line width (for stroking) |
| **Method** | `SetStrokeStyle(VGStrokeJoinStyle, VGStrokeCapStyle) error` | Set join and cap style  (for stroking) |
| **Method** | `SetStrokeDash(...float32) error` | Set the dash pattern (for stroking) |

The following system colors are defined:

| **Variable** | `khronos.VGColor` |
| -- | -- |
| `VGColorRed` | Primary Red |
| `VGColorGreen` | Primary Green |
| `VGColorBlue` | Primary Blue |
| `VGColorWhite` | Primary White |
| `VGColorBlack` | Primary Black |
| `VGColorPurple` | Purple |
| `VGColorCyan` | Cyan |
| `VGColorYellow` | Yellow |
| `VGColorDarkGrey` | Dark Grey |
| `VGColorLightGrey` | Light Grey |
| `VGColorMidGrey` | Mid Grey |

The following fill rules are defined:

| **Enum** | `khronos.VGFillRule` |
| -- | -- |
| `VG_STYLE_FILL_NONE` | Default fill rule (usually `VG_STYLE_FILL_EVENODD`)  |
| `VG_STYLE_FILL_NONZERO` | Non-zero fill rule |
| `VG_STYLE_FILL_EVENODD` | Even odd fill rule |

The following stroke cap styles are defined:

| **Enum** | `khronos.VGStrokeCapStyle` |
| -- | -- |
| `VG_STYLE_CAP_NONE` | Default cap style, or no change |
| `VG_STYLE_CAP_BUTT` | Butt end style, the default |
| `VG_STYLE_CAP_ROUND` | Round end style |
| `VG_STYLE_CAP_SQUARE` | Square end style |

The following stroke join styles are defined:

| **Enum** | `khronos.VGStrokeJoinStyle` |
| -- | -- |
| `VG_STYLE_JOIN_NONE` | Default join style, or no change |
| `VG_STYLE_JOIN_MITER` | Mitre join style, the default |
| `VG_STYLE_JOIN_ROUND` | Round join style |
| `VG_STYLE_JOIN_BEVEL` | Bevel join style |

The dash pattern for any strokes are defined using an on/off set of values.
For example:

TODO

## Creating and Drawing Paths

A path is a set of commands for drawing a shape. The `khronos.VGDriver` 
interface implements the following methods:

| **Interface** | `khronos.VGDriver` |
| -- | -- | -- |
| **Method** | `CreatePath() (VGPath, error)` | Create a path |
| **Method** | `DestroyPath(VGPath) error` | Destroy a path |

You can append drawing segments to the path using the following methods:

| **Interface** | `khronos.VGPath` |
| -- | -- | -- |
| **Method** | `MoveTo(VGPoint) error` | Set the beginning of the next segment to the point |
| **Method** | `LineTo(...VGPoint) error` | Draw one or more straight line segments |
| **Method** | `QuadTo(p1, p2 VGPoint) error` | Add a quadratic bezier segment from the last point, approaching control point p1, and ending at p2 |
| **Method** | `CubicTo(p1, p2, p3 VGPoint) error` | Add a cubic bezier egment from the last point, approaching control points p1 and p2, and ending at p3 |
| **Method** | `Close() error` | Close the current segment |

There are some utility functions to quickly create shapes:

| **Interface** | `khronos.VGPath` |
| -- | -- | -- |
| **Method** | `Line(start, end VGPoint) error` | Add a closed line to the path |
| **Method** | `Rect(origin, size VGPoint) error` | Add a close rectangle to the path |
| **Method** | `Ellipse(origin, diameter VGPoint) error` | Add a closed ellipse to the path |
| **Method** | `Circle(origin VGPoint, diameter float32) error` | Add a closed circle to the path |

Some other functions provide information about the path or affect the whole path:

| **Interface** | `khronos.VGPath` |
| -- | -- | -- |
| **Method** | `Clear() error` | Reset to an empty path, removing all path segments |
| **Method** | `Draw(stroke, fill VGPaint) error` | Stroke and fill path on current surface |
| **Method** | `Stroke(stroke VGPaint) error` | Stroke path on current surface |
| **Method** | `Fill(fill VGPaint) error` | Fill path on current surface |


# Links

  * (Wikipedia Entry](https://en.wikipedia.org/wiki/OpenVG)
  * [A Technical Introduction to OpenVG](https://www.khronos.org/assets/uploads/developers/library/siggraph2006/OpenKODE_Course/05_OpenVG-Overview.pdf)
  

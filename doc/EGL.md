
# Using EGL

The EGL layer roughly translates into the part of the system which creates
"surfaces" on the device display, and allows you to draw on them. The actual
drawing is performed by a higher-level API such as OpenVG or OpenGL.

In order to use the EGL layer, you'll need to create a "display" and then an
EGL object. For example:

```go
	// Create the display
    display, err := gopi.Open(rpi.DXDisplayConfig{
		Device:  this.Device,
		Display: config.Display,
	}, logger)
	if err != nil { /* handle error */ }
	defer display.Close()

	egl, err := gopi.Open(rpi.EGL{Display: display}, logger)
	if err != nil { /* handle error */ }
	defer egl.Close()
	
	/* do things here */
```

As shown in this example, once you've finished using your display, you will need
to `Close()` both the `egl` and `display` objects.

## Creating Surfaces

Once you have an EGL object, you can create surfaces on which you can draw. These
can either be a background (which covers the whole of the display screen) or
a normal surface (which can be positioned above the background).

When creating the surface, you'll need to specify the API you wish to "bind"
the drawing surface to. Typical values will be:

  * `"DX"` when you want to draw pixels onto the surface
  * `"OpenVG"` when you want to use 2D vector graphics
  * `"OpenGL"` when you want to use 3D vector graphics
  
Here's an example where you create a background and set the pixels in the
background to red. After every change of the bitmap, you'll need to "flush" the
surface to indiate the surface contents have changed:

```go
	bg, err := egl.CreateBackground("DX",1.0)
	if err != nil { /* handle error */ }
	defer egl.DestroySurface(bg)
	
	// Retrieve the bitmap resource for the surface
	bitmap, err := bg.GetBitmap()
	if err != nil { /* handle error */ }

	// Clear the background to red
	bitmap.Clear(khronos.EGLColorRed)
	egl.FlushSurface(bg)
```

Here's an example of creating a surface onto which you want to draw OpenVG
shapes. You need to provide the size of the surface in pixels and the origin
of the surface relative to the top left of the screen. The layer parameter
indicates where in the hierarchy of surfaces this one should be drawn, with
a higher number indicating it's the topmost surface.  The opacity parameter
can be between 0.0 and 1.0:

```go
    layer := uint16(2)
	opacity := float32(0.75)
	surface, err := egl.CreateSurface("OpenVG",khronos.EGLSize{ 200, 200 },khronos.EGLPoint{ 50, 50 },layer,opacity)
	if err != nil { /* handle error */ }
	defer egl.DestroySurface(surface)
```

As shown in these examples, you should call the `DestroySurface()` method when
you're done with the surfaces. This will remove the surface from the screen and
release any resources associated with it.

## Changing Layers and Opacity

Surfaces which are not the background surface can have their layer changed using
the following method:

TODO



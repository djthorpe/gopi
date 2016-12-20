
# Using fonts

You can render font faces either directly on EGL surfaces as bitmap fonts, or
using vector outlines on OpenVG surfaces. Before you can render text, you will
need to load fonts from either a file or from a folder which may contain many
font faces.

Create a font object using the following code:

```
  // Create a set of font faces
  fonts, err := gopi.Open(rpi.VGFont{
    PPI: 72
  }, logger)
  if err != nil { /* Handle error */ }
  defer fonts.Close()
```

When opening a `VGFont` instance, you should supply a PPI parameter which is
the density of pixels of your display in pixels per inch (PPI). If you set
this value to zero or omit it, then fonts will be sized in pixels rather than
points.

## Loading font faces from a file

You can then load fonts either individually or in a folder. To load a face from
a font file, use:

```
   face, err := fonts.OpenFace(path)
   if err != nil { /* Handle error */ }
   defer fonts.DestroyFace(face)
```

This loads the first face within the font file. There may be further faces within
the file. You can determine the number of faces within the file using the method
`GetNumFaces()`. For example,

```
   if face.GetNumFaces() > 0 {
      face, err := fonts.OpenFaceAtIndex(path,1)
      if err != nil { /* Handle error */ }   
      defer fonts.DestroyFace(face)
   }
```  

You can query the font families which have been loaded using the `GetFamilies()`
method and return a set of faces for one family using the `GetFaces()` method.
For example,

```
  for _, family := range fonts.GetFamilies() {
    for _, face := range fonts.GetFaces(family,khronos.VG_FONT_STYLE_ANY) {
	  // Do something with the face here
	}
  }
```

The `GetFaces()` method allows for one of the following flags:

  * khronos.VG_FONT_STYLE_ANY
  * khronos.VG_FONT_STYLE_BOLDITALIC
  * khronos.VG_FONT_STYLE_BOLD
  * khronos.VG_FONT_STYLE_ITALIC

## Loading font faces from a folder

You can use the method `OpenFacesAtPath` to load several font files at once,
for example from a folder. In order to use this, you need to implement a callback
function which can return true (to continue examining files and folders) or
false (to ignore a file or folder). For example,

```
  err := this.Fonts.OpenFacesAtPath(path, func(filename string, info os.FileInfo) bool {
    if strings.HasPrefix(info.Name(), ".") {
      // ignore hidden files and folders
      return false
    }
	if info.IsDir() {
      // recurse into folders
	  return true
	}
	if path.Ext(filename) == ".ttf" || path.Ext(filename) == ".TTF" {
	  // support TTF loading
	  return true
	}
	return false
  })
  if err != nil { /* handle error */ }
```

## Rendering text onto bitmap surfaces

You can render fonts directly onto bitmap surfaces. In order to do this,
use the `PaintText` method of a `khronos.EGLBitmap`. For example,

```
  surface := egl.CreateBackground("DX",1.0)
  surface.bitmap.PaintText("Hello, World!",face,khronos.EGLWhiteColor,surface.GetFrame().BottomLeft(),72.0)
```

The origin argument is the baseline of the text, so painting is both above and 
below the origin, where the text has characters with descenders (for example,
"p" and "q").

## Rendering text onto OpenVG surfaces

TODO

## Examples

There are two examples in the [examples](https://github.com/djthorpe/gopi/tree/master/examples/vgfont)
folder. 

  * `font_example.go` demonstrates how to load fonts and query which fonts are loaded.
  * `dx_example.go` demonstrates how to render text on bitmap surfaces.

For example,

```
  gopi% dx_example -fontpath etc/fonts -font Roboto -size 72
```





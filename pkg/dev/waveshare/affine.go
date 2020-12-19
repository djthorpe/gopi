package waveshare

import (
	"math"

	"golang.org/x/image/math/f64"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type AffineTransform f64.Aff3

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewAffineTransform() AffineTransform {
	return identityMatrix()
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (a AffineTransform) Scale(x, y float64) AffineTransform {
	return a.Multiply(scaleMatrix(x, y))
}

func (a AffineTransform) Translate(x, y float64) AffineTransform {
	return a.Multiply(translationMatrix(x, y))
}

func (a AffineTransform) Rotate(theta, x, y float64) AffineTransform {
	return a.Multiply(rotationMatrix(theta, x, y))
}

func (a AffineTransform) Multiply(b AffineTransform) AffineTransform {
	a[0] = a[0]*b[0] + a[1]*b[3]
	a[1] = a[0]*b[1] + a[1]*b[4]
	a[2] = a[0]*b[2] + a[1]*b[5] + a[2]
	a[3] = a[3]*b[0] + a[4]*b[3]
	a[4] = a[3]*b[1] + a[4]*b[4]
	a[5] = a[3]*b[2] + a[4]*b[5] + a[5]
	return a
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func identityMatrix() AffineTransform {
	return AffineTransform(f64.Aff3{
		1, 0, 0,
		0, 1, 0,
		// 0, 0, 1,
	})
}

func scaleMatrix(x, y float64) AffineTransform {
	r := NewAffineTransform()
	r[0] = x
	r[4] = y
	return r
}

func rotationMatrix(theta, x, y float64) AffineTransform {
	a := translationMatrix(-x, -y)
	r := identityMatrix()
	r[0] = math.Cos(theta)
	r[1] = -math.Sin(theta)
	r[3] = math.Sin(theta)
	r[4] = math.Cos(theta)
	return a.Multiply(r).Translate(x, y)
}

func translationMatrix(x, y float64) AffineTransform {
	b := NewAffineTransform()
	b[2] = x
	b[5] = y
	return b
}

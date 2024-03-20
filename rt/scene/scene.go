package scene

import (
	"image/color"
	"math"
)

type Sphere struct {
	X, Y, Z  float64
	Radius   float64
	Color    color.RGBA
	Specular float64
}

type Wall struct {
	X0, Y0, Z0, X1, Y1, Z1 float64
	Color                  color.RGBA
	Specular               float64
}

// Point light source
type Light struct {
	X, Y, Z float64
	Color   color.RGBA
}

// Wall intersect
func (w *Wall) Intersect(ox, oy, oz, dx, dy, dz float64) (hit bool, t float64) {
	txMin := (w.X0 - ox) / dx
	txMax := (w.X1 - ox) / dx
	tyMin := (w.Y0 - oy) / dy
	tyMax := (w.Y1 - oy) / dy
	tzMin := (w.Z0 - oz) / dz
	tzMax := (w.Z1 - oz) / dz

	tMin := math.Max(math.Max(math.Min(txMin, txMax), math.Min(tyMin, tyMax)), math.Min(tzMin, tzMax))
	tMax := math.Min(math.Min(math.Max(txMin, txMax), math.Max(tyMin, tyMax)), math.Max(tzMin, tzMax))

	if tMin > tMax || tMax <= 0 {
		return false, 0
	}

	return true, tMin
}

// Sphere intersect
func (s *Sphere) Intersect(ox, oy, oz, dx, dy, dz float64) (hit bool, t float64) {
	cx, cy, cz := s.X, s.Y, s.Z
	r := s.Radius

	lx := cx - ox
	ly := cy - oy
	lz := cz - oz

	dot := lx*dx + ly*dy + lz*dz
	if dot < 0 {
		return false, 0
	}

	d2 := lx*lx + ly*ly + lz*lz - dot*dot
	if d2 > r*r {
		return false, 0
	}

	thc := math.Sqrt(r*r - d2)
	t0 := dot - thc
	t1 := dot + thc

	if t0 < 0 {
		t = t1
	} else {
		t = t0
	}
	return true, t
}

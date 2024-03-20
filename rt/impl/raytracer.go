package impl

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"proj3/par"
	"proj3/scene"
	"sync"
)

func Parallel(
	wg *sync.WaitGroup, barrier *par.Barrier, me int, deQTrace, deQGauss *par.BDEQueue, traceQs, gaussQs []*par.BDEQueue, stealing bool,
	spheres []scene.Sphere, walls []scene.Wall, light scene.Light, fov, eyeX, eyeY, eyeZ float64,
	kernel [][]float64, dim int, img, result *image.RGBA, rays int) {

	defer wg.Done()
	// Ray tracing loop
	for {
		value, ok := deQTrace.PopBottom()
		if !ok {
			break
		}
		// Random work stealing
		if stealing {
			size := deQTrace.Size()
			if size >= 0 && rand.Intn(size+1) == size {
				victim := rand.Intn(len(traceQs))
				if victim != me {
					par.Balance(deQTrace, traceQs[victim])
				}
			}
		}
		x := value.X
		y := value.Y
		xf := float64(x)
		yf := float64(y)

		// Ray direction
		dx := math.Tan(fov/2) * ((xf - eyeX) / eyeX)
		dy := math.Tan(fov/2) * ((yf - eyeY) / eyeY)
		dz := 1.0
		len := math.Sqrt(dx*dx + dy*dy + dz*dz)
		dx /= len
		dy /= len
		dz /= len

		pixelColor := scene.TraceRay(eyeX, eyeY, eyeZ, dx, dy, dz, spheres, walls, light, rays)
		img.Set(x, y, pixelColor)
	}

	barrier.Signal()
	barrier.Wait()

	bounds := img.Bounds()
	// Gaussian Anti-aliasing loop
	for {
		value, ok := deQGauss.PopBottom()
		if !ok {
			break
		}
		// Random work stealing
		if stealing {
			size := deQGauss.Size()
			if size >= 0 && rand.Intn(size+1) == size {
				victim := rand.Intn(len(gaussQs))
				if victim != me {
					par.Balance(deQGauss, gaussQs[victim])
				}
			}
		}
		x := value.X
		y := value.Y
		var r, g, b float64
		_, _, _, a := img.At(x, y).RGBA()

		// Apply kernel
		for ky := 0; ky < dim; ky++ {
			for kx := 0; kx < dim; kx++ {
				px := x - 1 + kx
				py := y - 1 + ky
				if px < bounds.Min.X || py < bounds.Min.Y || px >= bounds.Max.X || py >= bounds.Max.Y {
					continue
				}
				pixel := img.At(px, py)
				pr, pg, pb, _ := pixel.RGBA()

				r += float64(pr) * kernel[kx][ky]
				g += float64(pg) * kernel[kx][ky]
				b += float64(pb) * kernel[kx][ky]
			}
		}

		result.Set(x, y, color.RGBA64{scene.Clamp(r), scene.Clamp(g), scene.Clamp(b), uint16(a)})
	}
}

func Sequential(
	spheres []scene.Sphere, walls []scene.Wall, light scene.Light, fov, eyeX, eyeY, eyeZ float64,
	kernel [][]float64, dim int, img, result *image.RGBA, rays int) {

	bounds := img.Bounds()
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			xf := float64(x)
			yf := float64(y)

			// Ray direction
			dx := math.Tan(fov/2) * ((xf - eyeX) / eyeX)
			dy := math.Tan(fov/2) * ((yf - eyeY) / eyeY)
			dz := 1.0
			len := math.Sqrt(dx*dx + dy*dy + dz*dz)
			dx /= len
			dy /= len
			dz /= len

			pixelColor := scene.TraceRay(eyeX, eyeY, eyeZ, dx, dy, dz, spheres, walls, light, rays)
			img.Set(x, y, pixelColor)
		}
	}

	// Gaussian Anti-aliasing loop
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			var r, g, b float64
			_, _, _, a := img.At(x, y).RGBA()

			// Apply kernel
			for ky := 0; ky < dim; ky++ {
				for kx := 0; kx < dim; kx++ {
					px := x - 1 + kx
					py := y - 1 + ky
					if px < bounds.Min.X || py < bounds.Min.Y || px >= bounds.Max.X || py >= bounds.Max.Y {
						continue
					}
					pixel := img.At(px, py)
					pr, pg, pb, _ := pixel.RGBA()

					r += float64(pr) * kernel[kx][ky]
					g += float64(pg) * kernel[kx][ky]
					b += float64(pb) * kernel[kx][ky]
				}
			}

			result.Set(x, y, color.RGBA64{scene.Clamp(r), scene.Clamp(g), scene.Clamp(b), uint16(a)})
		}
	}
}

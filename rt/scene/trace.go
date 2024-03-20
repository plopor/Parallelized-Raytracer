package scene

import (
	"image/color"
	"math"
)

func Reflect(ix, iy, iz, nx, ny, nz float64) (rx, ry, rz float64) {
	dot := ix*nx + iy*ny + iz*nz
	rx = ix - 2*dot*nx
	ry = iy - 2*dot*ny
	rz = iz - 2*dot*nz
	return rx, ry, rz
}

func PhongModel(light Light, spec, nx, ny, nz, eyeX, eyeY, eyeZ, hitX, hitY, hitZ float64) float64 {
	// Light vector
	lightDirX := light.X - hitX
	lightDirY := light.Y - hitY
	lightDirZ := light.Z - hitZ
	lightDist := math.Sqrt(lightDirX*lightDirX + lightDirY*lightDirY + lightDirZ*lightDirZ)
	lightDirX /= lightDist
	lightDirY /= lightDist
	lightDirZ /= lightDist

	// Camera vector
	vDirX := eyeX - hitX
	vDirY := eyeY - hitY
	vDirZ := eyeZ - hitZ
	vDist := math.Sqrt(vDirX*vDirX + vDirY*vDirY + vDirZ*vDirZ)
	vDirX /= vDist
	vDirY /= vDist
	vDirZ /= vDist

	i := 1.0
	// Ambient light
	ka := 0.2
	ambient := ka * i

	// Diffuse light
	dot := nx*lightDirX + ny*lightDirY + nz*lightDirZ
	kd := 0.4
	diffuse := kd * math.Max(0, dot) * i

	// Specular light
	rSpecX := -lightDirX + 2*dot*nx
	rSpecY := -lightDirY + 2*dot*ny
	rSpecZ := -lightDirZ + 2*dot*nz
	ks := 0.4
	specular := ks * math.Pow(math.Max(0, rSpecX*vDirX+rSpecY*vDirY+rSpecZ*vDirZ), spec) * i

	// Calculate final color using (simplified) Phong illumination model
	r := ambient + diffuse + specular

	return r
}

func TraceRay(eyeX, eyeY, eyeZ, dx, dy, dz float64, spheres []Sphere, walls []Wall, light Light, depth int) color.RGBA {
	if depth <= 0 {
		return color.RGBA{100, 100, 100, 255}
	}

	nearestSphere := -1
	minT := math.Inf(1)

	// Check intersection with spheres
	for i, sphere := range spheres {
		hit, dist := sphere.Intersect(eyeX, eyeY, eyeZ, dx, dy, dz)
		if hit && dist < minT {
			minT = dist
			nearestSphere = i
		}
	}

	nearestWall := -1
	minWallT := math.Inf(1)

	// Check intersection with walls
	for i, wall := range walls {
		hit, dist := wall.Intersect(eyeX, eyeY, eyeZ, dx, dy, dz)
		if hit && dist < minWallT {
			minWallT = dist
			nearestWall = i
		}
	}

	nearestObject := nearestSphere
	isSphere := true
	if minWallT < minT {
		nearestObject = nearestWall
		isSphere = false
	}

	if nearestObject != -1 {
		if isSphere {
			sphere := spheres[nearestObject]
			// Intersection
			hitX := eyeX + minT*dx
			hitY := eyeY + minT*dy
			hitZ := eyeZ + minT*dz

			// Normal
			nx := (hitX - sphere.X) / sphere.Radius
			ny := (hitY - sphere.Y) / sphere.Radius
			nz := (hitZ - sphere.Z) / sphere.Radius
			len := math.Sqrt(nx*nx + ny*ny + nz*nz)
			nx /= len
			ny /= len
			nz /= len

			r := PhongModel(light, sphere.Specular, nx, ny, nz, eyeX, eyeY, eyeZ, hitX, hitY, hitZ)
			reflectedDirX, reflectedDirY, reflectedDirZ := Reflect(dx, dy, dz, nx, ny, nz)
			reflectedColor := TraceRay(hitX, hitY, hitZ, reflectedDirX, reflectedDirY, reflectedDirZ, spheres, walls, light, depth-1)

			// Reflected light at 20% of total intensity
			return color.RGBA{
				uint8(float64(sphere.Color.R)*r*0.8 + float64(reflectedColor.R)*0.2),
				uint8(float64(sphere.Color.G)*r*0.8 + float64(reflectedColor.G)*0.2),
				uint8(float64(sphere.Color.B)*r*0.8 + float64(reflectedColor.B)*0.2),
				255,
			}
		} else {
			wall := walls[nearestObject]
			// Intersection
			hitX := eyeX + minWallT*dx
			hitY := eyeY + minWallT*dy
			hitZ := eyeZ + minWallT*dz

			// Normal
			vx := hitX - wall.X0
			vy := hitY - wall.Y0
			vz := hitZ - wall.Z0
			wx := wall.X1 - wall.X0
			wy := wall.Y1 - wall.Y0
			wz := wall.Z1 - wall.Z0

			nx := vy*wz - vz*wy
			ny := vz*wx - vx*wz
			nz := vx*wy - vy*wx
			len := math.Sqrt(nx*nx + ny*ny + nz*nz)
			nx /= len
			ny /= len
			nz /= len

			r := PhongModel(light, wall.Specular, nx, ny, nz, eyeX, eyeY, eyeZ, hitX, hitY, hitZ)
			reflectedDirX, reflectedDirY, reflectedDirZ := Reflect(dx, dy, dz, nx, ny, nz)
			reflectedColor := TraceRay(hitX, hitY, hitZ, reflectedDirX, reflectedDirY, reflectedDirZ, spheres, walls, light, depth-1)

			// Reflected light at 20% of total intensity
			return color.RGBA{
				uint8(float64(wall.Color.R)*r*0.8 + float64(reflectedColor.R)*0.2),
				uint8(float64(wall.Color.G)*r*0.8 + float64(reflectedColor.G)*0.2),
				uint8(float64(wall.Color.B)*r*0.8 + float64(reflectedColor.B)*0.2),
				255,
			}
		}
	}

	return color.RGBA{100, 100, 100, 255}
}

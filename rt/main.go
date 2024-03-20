package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"proj3/impl"
	"proj3/io"
	"proj3/par"
	"proj3/scene"
	"strconv"
	"sync"
	"time"
)

func main() {
	timeStart := time.Now()
	numThreads := 1
	stealing := false
	skew := false
	args := os.Args
	if len(args) > 1 {
		numThreads, _ = strconv.Atoi(args[1])
		if len(args) > 2 {
			stealing, _ = strconv.ParseBool(args[2])
			if len(args) > 3 {
				skew, _ = strconv.ParseBool(args[3])
			}
		}
	}

	width, height, depth, spheres, light, eyeZ, fov, dim, sigma, rays := io.ParseScene()
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Bounding box walls that scale with height, width, and depth of image
	walls := []scene.Wall{
		{X0: 0, Y0: 0, Z0: 0, X1: float64(width), Y1: 0, Z1: float64(depth), Color: color.RGBA{255, 0, 0, 255}, Specular: 0.5},
		{X0: 0, Y0: float64(height), Z0: 0, X1: float64(width), Y1: float64(height), Z1: float64(depth), Color: color.RGBA{0, 255, 0, 255}, Specular: 0.5},
		{X0: 0, Y0: 0, Z0: 0, X1: 0, Y1: float64(height), Z1: float64(depth), Color: color.RGBA{0, 0, 255, 255}, Specular: 0.5},
		{X0: float64(width), Y0: 0, Z0: 0, X1: float64(width), Y1: float64(height), Z1: float64(depth), Color: color.RGBA{255, 255, 0, 255}, Specular: 0.5},
		{X0: 0, Y0: 0, Z0: float64(depth), X1: float64(width), Y1: float64(height), Z1: float64(depth), Color: color.RGBA{255, 0, 255, 255}, Specular: 0.5},
	}
	eyeX, eyeY := float64(width/2), float64(height/2)

	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	kernel := scene.GaussianKernel(dim, sigma)

	if numThreads > 1 {
		var wg sync.WaitGroup
		barrier := par.NewBarrier(numThreads)

		var queue []par.Task
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				task := par.Task{X: x, Y: y}
				queue = append(queue, task)
			}
		}

		var traceQs []*par.BDEQueue
		var gaussQs []*par.BDEQueue

		// Split tasks evenly
		if !skew {
			step := int(math.Ceil(float64(width*height) / float64(numThreads)))
			for i, count := 0, 0; i < numThreads; i++ {
				subQueue := queue[count : count+step]
				deQTrace := par.NewBDEQueue(subQueue)
				deQGauss := par.NewBDEQueue(subQueue)
				traceQs = append(traceQs, deQTrace)
				gaussQs = append(gaussQs, deQGauss)
				count += step
			}
		} else {
			step := int(math.Ceil(float64(width*height) / float64(2*numThreads)))
			for i, count := 0, 0; i < numThreads; i++ {
				subQueue := queue[count : count+step]
				count += step
				if i >= numThreads/2 {
					subQueue = append(subQueue, queue[count:count+2*step]...)
					count = count + 2*step
				}
				//println(len(subQueue))
				deQTrace := par.NewBDEQueue(subQueue)
				deQGauss := par.NewBDEQueue(subQueue)
				traceQs = append(traceQs, deQTrace)
				gaussQs = append(gaussQs, deQGauss)
			}
		}

		for i := 0; i < numThreads; i++ {
			wg.Add(1)
			deQTrace := traceQs[i]
			deQGauss := gaussQs[i]
			go impl.Parallel(&wg, barrier, i, deQTrace, deQGauss, traceQs, gaussQs, stealing, spheres, walls, light, fov, eyeX, eyeY, eyeZ, kernel, dim, img, result, rays)
		}

		wg.Wait()
	} else {
		impl.Sequential(spheres, walls, light, fov, eyeX, eyeY, eyeZ, kernel, dim, img, result, rays)
	}

	file, err := os.Create("scene.png")
	if err != nil {
		println("Error creating output image:", err.Error())
		os.Exit(1)
	}
	defer file.Close()

	if err := png.Encode(file, result); err != nil {
		println("Error encoding scene:", err.Error())
		os.Exit(1)
	}

	fmt.Printf("%.2f\n", time.Since(timeStart).Seconds())
}

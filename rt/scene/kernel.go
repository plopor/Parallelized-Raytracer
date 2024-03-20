package scene

import "math"

func GaussianKernel(n int, s float64) [][]float64 {
	kernel := make([][]float64, n)
	sum := 0.0

	// Calculate kernel values
	for i := 0; i < n; i++ {
		kernel[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			x := float64(i - (n-1)/2)
			y := float64(j - (n-1)/2)
			kernel[i][j] = math.Exp(-(x*x+y*y)/(2.0*s*s)) / (2.0 * math.Pi * s * s)
			sum += kernel[i][j]
		}
	}

	// Normalize
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			kernel[i][j] /= sum
		}
	}

	return kernel
}

func Clamp(comp float64) uint16 {
	return uint16(math.Min(65535, math.Max(0, comp)))
}

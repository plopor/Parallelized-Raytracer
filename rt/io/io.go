package io

import (
	"encoding/json"
	"fmt"
	"os"
	"proj3/scene"
)

type Scene struct {
	Z      float64 `json:"Z"`
	FOV    float64 `json:"FOV"`
	Width  int     `json:"Width"`
	Height int     `json:"Height"`
	Depth  float64 `json:"Depth"`
	Kernel int     `json:"Kernel"`
	Sigma  float64 `json:"Sigma"`
	Rays   int     `json:"Rays"`
}

type Data struct {
	Spheres []scene.Sphere `json:"spheres"`
	Light   scene.Light    `json:"light"`
	Scene   Scene          `json:"scene"`
}

func ParseScene() (int, int, float64, []scene.Sphere, scene.Light, float64, float64, int, float64, int) {
	file, err := os.Open("scene.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	var data Data
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding scene:", err)
	}

	return data.Scene.Width, data.Scene.Height, data.Scene.Depth, data.Spheres, data.Light, data.Scene.Z, data.Scene.FOV, data.Scene.Kernel, data.Scene.Sigma, data.Scene.Rays
}

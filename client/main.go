package main

import (
	"math"
//	"math/rand"
//	"log"
//	"sync"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
)

func main() {

	// Create application and scene
	a := app.App()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	pDim := float32(100)
	r := float32(1)

	tileMap := geometry.NewPlane(pDim, pDim)
	tileMapMaterial := material.NewStandard(math32.NewColor("DarkBlue"))
	tileMapMesh := graphic.NewMesh(tileMap, tileMapMaterial)
	scene.Add(tileMapMesh)

	tickLen := 30 * time.Millisecond
	var tickMux sync.RWMutex
	var tick time.Time
	go func() {
		for {
			tickMux.Lock()
			tick = time.Now()
			tickMux.Unlock()
		}
	}()

	// Render the entire world. Bad idea.
	// Unity has hard-capped 5k entity limit (see http://answers.unity.com/answers/408712/view.html).
	//
	// We need to cull (for occlusion) and delete rendered objects when outside FOV.
	for i := 0; i < 10000; i++ {
		tank := geometry.NewSphereSector(float64(r), 2, 2, 0, math.Pi, 0, math.Pi)
		tankMaterial := material.NewStandard(math32.NewColor("Red"))
		tankMesh := graphic.NewMesh(tank, tankMaterial)
		scene.Add(tankMesh)

		go func() {
			x := (pDim - 2 * r) * rand.Float32() - (pDim / float32(2) - r)
			y := (pDim - 2 * r) * rand.Float32() - (pDim / float32(2) - r)
			for {
				tickMux.RLock()
				s := tick
				tickMux.RUnlock()

				dx := 0.1 * (rand.Float32() - 0.5)
				dy := 0.1 * (rand.Float32() - 0.5)
				if x + dx > 5 || x + dx < -5 {
					dx = 0
				}
				if y + dy > 5 || y + dy < -5 {
					dy = 0
				}
				d := math32.NewVector3(x + dx, y + dy, 0)
				tankMesh.SetPositionVec(d)

				timeToSleep := tickLen - time.Now().Sub(s)
				if timeToSleep < 0 {
					log.Printf("too SLOW %v", timeToSleep)
				}
				time.Sleep(timeToSleep)
			}
		}()
	}

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}

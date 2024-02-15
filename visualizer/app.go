package visualizer

import (
	"eve-firmware/arm"
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

var (
	joints []*graphic.Mesh
	scene  *core.Node
)

func ReloadPos() {
	for i := 0; i < 7; i++ {
		mesh := joints[i]
		matrix := math32.NewMatrix4()
		arm.POSITION.HTMatrices[i].Print()
		t := arm.HTMFromTo(0, i+1).D()
		matrix.Set(float32(t[0][0]), float32(t[0][1]), float32(t[0][2]), float32(t[0][3]),
			float32(t[1][0]), float32(t[1][1]), float32(t[1][2]), float32(t[1][3]),
			float32(t[2][0]), float32(t[2][1]), float32(t[2][2]), float32(t[2][3]),
			float32(t[3][0]), float32(t[3][1]), float32(t[3][2]), float32(t[3][3]))
		mesh.SetMatrix(matrix)
	}
}

func RunApp() {
	// Create application and scene
	a := app.App()
	scene = core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 100)
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

	// Create a blue torus and add it to the scene
	geom := geometry.NewCylinder(20, 50, 32, 10, true, true)
	base := material.NewStandard(math32.NewColor("Yellow"))
	mat := material.NewStandard(math32.NewColor("DarkBlue"))

	for i := 0; i < 7; i++ {
		mesh := graphic.NewMesh(geom, mat)
		joints = append(joints, mesh)
		scene.Add(mesh)
		mesh.SetPosition(0, 0, 0)
	}
	ReloadPos()

	joints[0].SetMaterial(base)

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(5))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}
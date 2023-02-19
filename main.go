package main

import (
	"github.com/UpsilonDiesBackwards/behngine_epq/camera"
	"github.com/UpsilonDiesBackwards/behngine_epq/input"
	"github.com/UpsilonDiesBackwards/behngine_epq/primitives"
	"github.com/UpsilonDiesBackwards/behngine_epq/windowing"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
	"runtime"
	"time"
)

// Geometry
var (
	triangle = []float32{
		0, 0.5, 0, // top
		-0.5, -0.5, 0, // left
		0.5, -0.5, 0, // right
	}

	square = []float32{
		// first triangle
		0.5, 0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, 0.5, 0.0,

		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
		-0.5, 0.5, 0.0,

		// second triangle
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,

		-0.5, 0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
	}
	square_i = []uint32{
		0, 1, 2, // first triangle of square
		2, 3, 0, // second triangle of square
	}
)
var DeltaTime float64

var uI = &input.UserInput{} // Create pointer to UserInput struct

// FPS Counter
var startTime = time.Now()
var frameCount int
var FPS float64

var object = &primitives.ObjectPrimitive{}

func main() {
	runtime.LockOSThread()
	defer glfw.Terminate()

	appWindow := windowing.NewWindow(1024, 768, "3D Rendering Engine")

	//object.objectV, object.objectN, object.objectTN, object.objectI = primitives.CreateNewOBJ("primitives/tower.obj")
	//primitives.CreateOBJObject("primitives/cube.obj", "Test", appWindow, &camera.Camera_Viewport)

	var previousTime = time.Now()

	// Primary program loop
	for !appWindow.ShouldClose() {
		// Calculate delta time
		currentTime := time.Now()
		DeltaTime = currentTime.Sub(previousTime).Seconds()
		previousTime = currentTime

		appWindow.Execute(&camera.Camera_Viewport)

		// For each frame, check for user input
		err := ProgramInputLoop(appWindow, DeltaTime, &camera.Camera_Viewport, uI)
		if err != nil {
			log.Fatalln(err)
		}

		// Handle FPS counter
		frameCount++
		if time.Since(startTime) >= time.Second {
			// Calculate the FPS
			FPS = float64(frameCount) / time.Since(startTime).Seconds()

			// Print the FPS
			//str := strconv.FormatFloat(FPS, 'f', 1, 64)
			//fmt.Println("FPS: ", str)

			// Reset the frame count and start time
			frameCount = 0
			startTime = time.Now()
		}

		//// Update the window content
		//drawWindowContent(VAO, program)
	}
}

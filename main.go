package main

import (
	"fmt"
	"github.com/UpsilonDiesBackwards/behngine_epq/camera"
	"github.com/UpsilonDiesBackwards/behngine_epq/input"
	"github.com/UpsilonDiesBackwards/behngine_epq/objLoader"
	"github.com/UpsilonDiesBackwards/behngine_epq/shaders"
	"github.com/UpsilonDiesBackwards/behngine_epq/windowing"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"runtime"
	"time"
)

type PerspectiveBlock struct {
	project mgl32.Mat4
	camera  mgl32.Mat4
	model   mgl32.Mat4
}

var UBO uint32

var fov = float32(60.0)
var projectionTransform = mgl32.Perspective(mgl32.DegToRad(fov),
	float32(800/600), 0.1, 2000)

var DeltaTime float64

var userInput = &input.UserInput{} // Create pointer to UserInput struct

func main() {
	runtime.LockOSThread()
	defer glfw.Terminate()

	appWindow, program := CreateWindow(1024, 768, "3D Rendering Engine")

	objLoader.CreateObject("cube.obj", "CUBE", mgl32.Vec3{-4, 0, 0}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 1, 1})
	objLoader.CreateObject("cube.obj", "CUBE", mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 1, 1})
	objLoader.CreateObject("cube.obj", "CUBE", mgl32.Vec3{4, 0, 0}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 1, 1})

	var previousTime = time.Now()

	// Primary program loop
	for !windowing.ShouldClose {
		// Calculate delta time
		currentTime := time.Now()
		DeltaTime = currentTime.Sub(previousTime).Seconds()
		previousTime = currentTime

		// For each frame, check for user input
		err := ProgramInputLoop(appWindow, DeltaTime, &camera.Camera_Viewport, userInput)
		if err != nil {
			log.Fatalln(err)
		}

		drawWindowContent(program)

		windowing.EnableFPSCounter(DeltaTime)

		// Update the window
		appWindow.SwapBuffers()
		glfw.PollEvents()
	}
}

func CreateWindow(width, height int, title string) (*glfw.Window, uint32) {
	fmt.Printf("Initializing GLFW... ")
	if err := glfw.Init(); err != nil {
		panic(err)
	} else {
		fmt.Println("Success!")
	}

	// Establish window hints
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Create a new glfw window called `appWindow`
	appWindow, err := glfw.CreateWindow(
		width,
		height,
		title,
		nil, nil)
	if err != nil {
		panic(err)
	}
	appWindow.MakeContextCurrent()

	appWindow.SetKeyCallback(userInput.KeyCallback)
	appWindow.SetCursorPosCallback(userInput.MouseCallBack)

	fmt.Printf("Initializing OpenGL... ")
	if err := gl.Init(); err != nil {
		panic(err)
	} else {
		fmt.Println("Success!")
	}

	_GLVersion := gl.GoStr(gl.GetString(gl.VERSION)) // Get current OpenGL version
	log.Println("Using OpenGL version:", _GLVersion)

	// Create vShader and fShader
	vShader, err := shaders.ShaderCompiler(shaders.SHADER_DEFAULT_V, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fShader, err := shaders.ShaderCompiler(shaders.SHADER_DEFAULT_F, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	// Create a new shader program, then attach shaders, then link and use it.
	glProgram := gl.CreateProgram()
	gl.AttachShader(glProgram, vShader)
	gl.AttachShader(glProgram, fShader)
	gl.LinkProgram(glProgram)
	gl.UseProgram(glProgram)

	return appWindow, glProgram
}

func drawWindowContent(program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // Clear the color and depth buffer bits
	gl.ClearColor(0.52, 0.80, 0.96, 1.0)

	gl.UseProgram(program)
	if err := gl.GetError(); err != 0 {
		log.Printf("error using program: %v\n", err)
	}

	// Loop through all objects and render them
	for _, obj := range objLoader.Objects {
		gl.BindVertexArray(obj.VAO)
		if err := gl.GetError(); err != 0 {
			log.Printf("error binding vertex array object for object: %v,  reason: &v\n", obj, err)
		}

		obj.CreateModelMatrix()

		block := createPerspectiveBlock(projectionTransform, ViewportTransform, obj.ModelMatrix)
		UBO := createUBO(&block.project, &block.camera, &block.model)

		gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, UBO)
		if err := gl.GetError(); err != 0 {
			log.Printf("error binding UBO: %v\n", err)
		}

		gl.DrawElements(gl.TRIANGLES, int32(len(obj.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	}

	gl.Enable(gl.DEPTH_TEST)
}

func createUBO(project, camera, model *mgl32.Mat4) uint32 {
	var ubo uint32
	gl.GenBuffers(1, &ubo)
	gl.BindBuffer(gl.UNIFORM_BUFFER, ubo)
	gl.BufferData(gl.UNIFORM_BUFFER, 3*16*4, nil, gl.STREAM_DRAW)
	gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, ubo)

	gl.BindBuffer(gl.UNIFORM_BUFFER, ubo)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, 16*4, gl.Ptr(&project[0]))
	gl.BufferSubData(gl.UNIFORM_BUFFER, 16*4, 16*4, gl.Ptr(&camera[0]))
	gl.BufferSubData(gl.UNIFORM_BUFFER, 32*4, 16*4, gl.Ptr(&model[0]))

	return ubo
}

func createPerspectiveBlock(projection, viewport, world mgl32.Mat4) PerspectiveBlock {
	return PerspectiveBlock{
		project: projection,
		camera:  viewport,
		model:   world,
	}
}

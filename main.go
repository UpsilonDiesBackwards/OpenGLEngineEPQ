package main

import (
	"fmt"
	"github.com/UpsilonDiesBackwards/behngine_epq/camera"
	"github.com/UpsilonDiesBackwards/behngine_epq/input"
	"github.com/UpsilonDiesBackwards/behngine_epq/shaders"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"runtime"
	"time"
)

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

type PerspectiveBlock struct {
	project *mgl32.Mat4
	camera  *mgl32.Mat4
	world   *mgl32.Mat4
}

var fov = float32(60.0)
var projectionTransform = mgl32.Perspective(mgl32.DegToRad(fov),
	float32(800/600), 0.1, 2000)

var UBO uint32
var DeltaTime float64

var uI = &input.UserInput{} // Create pointer to UserInput struct

func main() {
	runtime.LockOSThread()

	appWindow := initGLFW()
	defer glfw.Terminate()

	translate := mgl32.Vec3{0, 0, -3} // Translation vector

	angle := float32(1 / 4) // Rotation angle (in radians) and axis
	axis := mgl32.Vec3{0, 1, 0}

	scale := mgl32.Vec3{2, 2, 2} // Scale factors

	// Create translation, rotation, and scale matrices
	translateMatrix := mgl32.Translate3D(translate.X(), translate.Y(), translate.Z())
	rotateMatrix := mgl32.HomogRotate3D(angle, axis)
	scaleMatrix := mgl32.Scale3D(scale.X(), scale.Y(), scale.Z())

	// Multiply the matrices to create the world matrix
	world := translateMatrix.Mul4(rotateMatrix.Mul4(scaleMatrix))

	var previousTime = time.Now()

	program := initGL()
	VAO := CreateVAO(square)
	if VAO == 0 {
		log.Fatalln("Error creating VAO")
	}

	// Primary program loop
	for !appWindow.ShouldClose() {
		// Measure the time since the last frame
		currentTime := time.Now()
		DeltaTime = currentTime.Sub(previousTime).Seconds()
		previousTime = currentTime

		// For each frame, check for user input
		err := ProgramInputLoop(appWindow, DeltaTime, &camera.Camera_Viewport, uI)
		if err != nil {
			log.Fatalln(err)
		}

		// Update the window content
		drawWindowContent(VAO, appWindow, program)

		block := PerspectiveBlock{
			project: &projectionTransform,
			camera:  &ViewportTransform,
			world:   &world,
		}
		UBO = CreateUBO(block)

		// Update the window
		appWindow.SwapBuffers()
		glfw.PollEvents()
	}
}

func initGLFW() *glfw.Window {
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
		1024,
		768,
		"3D Rendering Engine",
		nil, nil)
	if err != nil {
		panic(err)
	}
	appWindow.MakeContextCurrent()
	appWindow.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	return appWindow
}

func initGL() uint32 {
	fmt.Printf("Initializing OpenGL... ")
	if err := gl.Init(); err != nil {
		panic(err)
	} else {
		fmt.Println("Success!")
	}

	_GLVersion := gl.GoStr(gl.GetString(gl.VERSION)) // Get current OpenGL version
	log.Println("Using OpenGL version:", _GLVersion)

	// Create vShader and fShader
	vShader, err := shaders.CompilerCompiler(shaders.SHADER_DEFAULT_V, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fShader, err := shaders.CompilerCompiler(shaders.SHADER_DEFAULT_F, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	// Create a new shader program, then attach shaders, then link and use it.
	glProgram := gl.CreateProgram()
	gl.AttachShader(glProgram, vShader)
	gl.AttachShader(glProgram, fShader)
	gl.LinkProgram(glProgram)
	gl.UseProgram(glProgram)
	return glProgram

	gl.UseProgram(glProgram)
	return glProgram
}

func drawWindowContent(VAO uint32, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // Clear the color and depth buffer bits
	gl.ClearColor(0.52, 0.80, 0.96, 1.0)

	gl.BindVertexArray(VAO)
	if err := gl.GetError(); err != 0 {
		log.Printf("error binding vertex array object: %v\n", err)
	}

	gl.UseProgram(program)
	if err := gl.GetError(); err != 0 {
		log.Printf("error using program: %v\n", err)
	}

	gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, UBO)
	if err := gl.GetError(); err != 0 {
		log.Printf("error binding UBO: %v\n", err)
	}

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))

	gl.Enable(gl.DEPTH_TEST)
}

func CreateVAO(vertices []float32) uint32 {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	//gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return VAO
}

func CreateUBO(block PerspectiveBlock) uint32 {
	var ubo uint32
	gl.GenBuffers(1, &ubo)
	gl.BindBuffer(gl.UNIFORM_BUFFER, ubo)
	gl.BufferData(gl.UNIFORM_BUFFER, 3*16*4, nil, gl.STREAM_DRAW)
	gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, ubo)

	gl.BindBuffer(gl.UNIFORM_BUFFER, ubo)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, 16*4, gl.Ptr(&block.project[0]))
	gl.BufferSubData(gl.UNIFORM_BUFFER, 16*4, 16*4, gl.Ptr(&block.camera[0]))
	gl.BufferSubData(gl.UNIFORM_BUFFER, 32*4, 16*4, gl.Ptr(&block.world[0]))

	return ubo
}

func CheckGLErrors() {
	glerror := gl.GetError()
	if glerror == gl.NO_ERROR {
		return
	}

	fmt.Printf("gl.GetError() reports")
	for glerror != gl.NO_ERROR {
		fmt.Printf(" ")
		switch glerror {
		case gl.INVALID_ENUM:
			fmt.Printf("GL_INVALID_ENUM")
		case gl.INVALID_VALUE:
			fmt.Printf("GL_INVALID_VALUE")
		case gl.INVALID_OPERATION:
			fmt.Printf("GL_INVALID_OPERATION")
		case gl.STACK_OVERFLOW:
			fmt.Printf("GL_STACK_OVERFLOW")
		case gl.STACK_UNDERFLOW:
			fmt.Printf("GL_STACK_UNDERFLOW")
		case gl.OUT_OF_MEMORY:
			fmt.Printf("GL_OUT_OF_MEMORY")
		default:
			fmt.Printf("%d", glerror)
		}
		glerror = gl.GetError()
	}
	fmt.Printf("\n")
}

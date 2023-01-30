package main

import (
	"fmt"
	"github.com/UpsilonDiesBackwards/behngine_epq/camera"
	"github.com/UpsilonDiesBackwards/behngine_epq/input"
	"github.com/UpsilonDiesBackwards/behngine_epq/primitives"
	"github.com/UpsilonDiesBackwards/behngine_epq/shaders"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
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

// FPS Counter
var startTime = time.Now()
var frameCount int
var FPS float64

// Object
type ObjectPrimitive struct {
	objectV []float32
	objectI []uint32
}

var object = &ObjectPrimitive{}

func main() {
	runtime.LockOSThread()
	defer glfw.Terminate()

	appWindow := initGLFW()

	world := createWorldMatrix() // Create world matrix
	program := initGL()

	object.objectV, object.objectI = primitives.CreateNewOBJ("primitives/cube.obj")

	fmt.Println(object.objectI)

	VAO := createVAO(object.objectV, object.objectI)
	if VAO == 0 {
		log.Fatalln("Error creating VAO")
	}
	var previousTime = time.Now()

	// Primary program loop
	for !appWindow.ShouldClose() {
		// Calculate delta time
		currentTime := time.Now()
		DeltaTime = currentTime.Sub(previousTime).Seconds()
		previousTime = currentTime

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

		// Update the window content
		drawWindowContent(VAO, program)

		block := createPerspectiveBlock(&projectionTransform, &ViewportTransform, &world)
		UBO = createUBO(block)

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
	return glProgram
}

func drawWindowContent(VAO uint32, program uint32) {
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

	fmt.Println(gl.GetError())

	//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
	gl.DrawElements(gl.TRIANGLES, int32(len(object.objectI)), gl.UNSIGNED_INT, nil)

	gl.Enable(gl.DEPTH_TEST)
}

func createVAO(vertices []float32, indices []uint32) uint32 {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	var EBO uint32
	gl.GenBuffers(1, &EBO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	//gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return VAO
}

func createUBO(block PerspectiveBlock) uint32 {
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

func createWorldMatrix() mgl32.Mat4 {
	translate := mgl32.Vec3{0, 0, -3}
	angle := float32(1 / 4)
	axis := mgl32.Vec3{0, 1, 0}
	scale := mgl32.Vec3{2, 2, 2}

	translateMatrix := mgl32.Translate3D(translate.X(), translate.Y(), translate.Z())
	rotateMatrix := mgl32.HomogRotate3D(angle, axis)
	scaleMatrix := mgl32.Scale3D(scale.X(), scale.Y(), scale.Z())

	return translateMatrix.Mul4(rotateMatrix.Mul4(scaleMatrix))
}

func createPerspectiveBlock(projection *mgl32.Mat4, viewport *mgl32.Mat4, world *mgl32.Mat4) PerspectiveBlock {
	return PerspectiveBlock{
		project: projection,
		camera:  viewport,
		world:   world,
	}
}

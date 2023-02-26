package main

import (
	"fmt"
	"github.com/UpsilonDiesBackwards/behngine_epq/camera"
	"github.com/UpsilonDiesBackwards/behngine_epq/entities"
	"github.com/UpsilonDiesBackwards/behngine_epq/input"
	"github.com/UpsilonDiesBackwards/behngine_epq/shaders"
	"github.com/UpsilonDiesBackwards/behngine_epq/windowing"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"runtime"
	"time"
	"unsafe"
)

type PerspectiveBlock struct {
	project mgl32.Mat4
	camera  mgl32.Mat4
	model   mgl32.Mat4
}

type LightingBlock struct {
	transform mgl32.Mat4
	camera    mgl32.Mat4
	model     mgl32.Mat4
}

var fov = float32(60.0)
var projectionTransform = mgl32.Perspective(mgl32.DegToRad(fov),
	float32(800/600), 0.1, 2000)

var DeltaTime float64

var userInput = &input.UserInput{} // Create pointer to UserInput struct

func main() {
	runtime.LockOSThread()
	defer glfw.Terminate()

	appWindow, program, lightProgram := CreateWindow(1024, 768, "3D Rendering Engine")

	entities.CreateObject("tree.obj", "CUBE", mgl32.Vec3{-4, 0, 0}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 1, 1})
	entities.CreateObject("cube.obj", "CUBE", mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 1, 1})
	entities.CreateObject("tree.obj", "CUBE", mgl32.Vec3{4, 0, 0}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 1, 1})

	entities.CreateLight(mgl32.Vec3{-2, 0, 1}, mgl32.Vec3{1, 0.5, 0.5})

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

		drawWindowContent(program, lightProgram)

		windowing.EnableFPSCounter(DeltaTime)

		// Update the window
		appWindow.SwapBuffers()
		glfw.PollEvents()
	}
}

func CreateWindow(width, height int, title string) (*glfw.Window, uint32, uint32) {
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
	vShader, err := shaders.ShaderCompiler(shaders.SHADER_PHONG_V, gl.VERTEX_SHADER)
	if err != nil {
		fmt.Println("Error compiling vertex shader: ", err)
	}

	fShader, err := shaders.ShaderCompiler(shaders.SHADER_PHONG_F, gl.FRAGMENT_SHADER)
	if err != nil {
		fmt.Println("Error compiling fragment shader: ", err)
	}

	oProgram := gl.CreateProgram()
	gl.AttachShader(oProgram, vShader)
	gl.AttachShader(oProgram, fShader)
	gl.LinkProgram(oProgram)

	lShader, err := shaders.ShaderCompiler(shaders.SHADER_LIGHT_F, gl.FRAGMENT_SHADER)
	if err != nil {
		fmt.Println("Error compiling light shader: ", err)
	}

	lProgram := gl.CreateProgram()
	gl.AttachShader(lProgram, vShader)
	gl.AttachShader(lProgram, lShader)
	gl.LinkProgram(lProgram)

	return appWindow, oProgram, lProgram
}

func drawWindowContent(objectProgram, lightProgram uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // Clear the color and depth buffer bits
	gl.ClearColor(0.52, 0.80, 0.96, 1.0)

	// Use object Program
	gl.UseProgram(objectProgram)
	if err := gl.GetError(); err != 0 {
		log.Printf("error using object Program: %v\n", err)
	}

	// Loop through all objects
	for _, obj := range entities.Objects {
		gl.BindVertexArray(obj.VAO)
		if err := gl.GetError(); err != 0 {
			log.Printf("error binding vertex array object for object: %v,  reason: %v\n", obj.VAO, err)
		}

		objectColorLoc := gl.GetUniformLocation(objectProgram, gl.Str("objectColor\x00"))
		lightColorLoc := gl.GetUniformLocation(objectProgram, gl.Str("lightColor\x00"))
		lightPosLoc := gl.GetUniformLocation(objectProgram, gl.Str("lightPos\x00"))
		gl.Uniform3f(lightPosLoc, 1.2, 1.0, 2.0)     // Example light position
		gl.Uniform3f(objectColorLoc, 1.0, 0.5, 0.31) // Example object color
		gl.Uniform3f(lightColorLoc, 1.0, 1.0, 1.0)   // Example light color

		// Set perspective block
		obj.CreateModelMatrix()
		block := createPerspectiveBlock(projectionTransform, ViewportTransform, obj.ModelMatrix)
		UBO := createPerspectiveUBO(&block.project, &block.camera, &block.model)

		// Bind and draw objects
		gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, UBO)
		if err := gl.GetError(); err != 0 {
			log.Printf("error binding object uniform buffer base: %v\n", err)
		}
		gl.DrawElements(gl.TRIANGLES, int32(len(obj.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	}
	// unbind object
	gl.BindVertexArray(0)
	gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, 0)

	lightPos := mgl32.Vec3{0.6, 1, 0.1}
	lightTransform := mgl32.Translate3D(lightPos.X(), lightPos.Y(), lightPos.Z()).Mul4(
		mgl32.Scale3D(0.2, 0.2, 0.2))

	// Use light program
	gl.UseProgram(lightProgram)
	if err := gl.GetError(); err != 0 {
		log.Printf("error using light Program: %v\n", err)
	}

	// Loop through all lights
	for _, l := range entities.Lights {
		gl.BindVertexArray(l.VAO)
		if err := gl.GetError(); err != 0 {
			log.Printf("error binding vertex array object for light: %v,  reason: %v\n", l.VAO, err)
		}

		// Create lighting block
		l.CreateLightModelMatrix()
		block := createLightingBlock(lightTransform, ViewportTransform, l.ModelMatrix)
		UBO := createLightingUBO(&block.transform, &block.camera, &block.model, lightProgram)

		// Bind and draw
		gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, UBO)
		if err := gl.GetError(); err != 0 {
			log.Printf("error binding light uniform buffer base: %v for buffer %b\n", err, UBO)
		}

		gl.BufferSubData(gl.UNIFORM_BUFFER, 0, int(unsafe.Sizeof(block)), unsafe.Pointer(&block))
		gl.DrawElements(gl.TRIANGLES, int32(len(l.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	}
	// unbind light
	gl.BindVertexArray(0)
	gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, 0)

	gl.Enable(gl.DEPTH_TEST)
}

func createPerspectiveUBO(project, camera, model *mgl32.Mat4) uint32 {
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

func createLightingUBO(transform, camera, model *mgl32.Mat4, program uint32) uint32 {
	var ubo uint32
	gl.GenBuffers(1, &ubo)
	gl.BindBuffer(gl.UNIFORM_BUFFER, ubo)
	gl.BufferData(gl.UNIFORM_BUFFER, 3*16*4, nil, gl.STREAM_DRAW)
	gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, ubo)

	gl.BindBuffer(gl.UNIFORM_BUFFER, ubo)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, 16*4, gl.Ptr(&transform[0]))
	gl.BufferSubData(gl.UNIFORM_BUFFER, 16*4, 16*4, gl.Ptr(&camera[0]))
	gl.BufferSubData(gl.UNIFORM_BUFFER, 32*4, 16*4, gl.Ptr(&model[0]))

	// Bind the buffer to the lightBlock struct
	lightBlockBindingIndex := uint32(1) // This should match the binding index used in the shader
	lightBlockSize := int32(3 * 16 * 4) // Size of the lighting block in bytes
	gl.BindBufferBase(gl.UNIFORM_BUFFER, lightBlockBindingIndex, ubo)
	lightBlockOffset := int32(0)
	//gl.UniformBlockBinding(program, gl.GetUniformBlockIndex(program, gl.Str("LightingBlock\x00")), lightBlockBindingIndex)
	gl.BindBufferRange(gl.UNIFORM_BUFFER, lightBlockBindingIndex, ubo, int(lightBlockOffset), int(lightBlockSize))

	return ubo
}

func createPerspectiveBlock(projection, viewport, world mgl32.Mat4) PerspectiveBlock {
	return PerspectiveBlock{
		project: projection,
		camera:  viewport,
		model:   world,
	}
}

func createLightingBlock(position, viewport, world mgl32.Mat4) LightingBlock {
	return LightingBlock{
		transform: position,
		camera:    viewport,
		model:     world,
	}
}

func printAttachedShaders(program uint32) {
	shaders := make([]uint32, 2) // Initialize with a capacity of 2
	count := int32(len(shaders))
	gl.GetAttachedShaders(program, count, nil, &shaders[0])

	fmt.Printf("Program %d has %d shaders attached:\n", program, count)
	for _, shader := range shaders {
		var shaderType int32
		gl.GetShaderiv(shader, gl.SHADER_TYPE, &shaderType)
		fmt.Printf("\tShader %d: type %d\n", shader, shaderType)
	}
}

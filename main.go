package main

import (
	"fmt"
	"github.com/UpsilonDiesBackwards/behngine_epq/shaders"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
	"runtime"
)

var s = shaders.Program{}

var (
	triangle = []float32{
		0, 0.5, 0, // top
		-0.5, -0.5, 0, // left
		0.5, -0.5, 0, // right
	}
)

func main() {
	runtime.LockOSThread()

	appWindow := initGLFW()
	defer glfw.Terminate()

	program := initGL()

	VAO := createVAO(triangle)

	for !appWindow.ShouldClose() {
		drawWindowContent(VAO, appWindow, program)
	}
}

func initGLFW() *glfw.Window {
	fmt.Printf("Initializing GLFW... ")
	if err := glfw.Init(); err != nil {
		panic(err)
	} else {
		fmt.Println("Success!")
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	appWindow, err := glfw.CreateWindow(
		800,
		600,
		"3D Rendering Engine",
		nil, nil)
	if err != nil {
		panic(err)
	}
	appWindow.MakeContextCurrent()

	return appWindow
}

func initGL() uint32 {
	fmt.Printf("Initializing OpenGL... ")
	if err := gl.Init(); err != nil {
		panic(err)
	} else {
		fmt.Println("Success!")
	}

	_GLVersion := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("Using OpenGL version:", _GLVersion)

	vShader, err := shaders.Shader_compiler(shaders.SHADER_DEFAULT_V, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fShader, err := shaders.Shader_compiler(shaders.SHADER_DEFAULT_F, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	glProgram := gl.CreateProgram()
	gl.AttachShader(glProgram, vShader)
	gl.AttachShader(glProgram, fShader)
	gl.LinkProgram(glProgram)
	gl.UseProgram(glProgram)

	return glProgram
}

func drawWindowContent(VAO uint32, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0.52, 0.80, 0.96, 1.0)

	gl.BindVertexArray(VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))
	gl.Enable(gl.DEPTH_TEST)

	glfw.PollEvents()
	window.SwapBuffers()
}

func createVAO(vertices []float32) uint32 {
	// Generate the vertex buffer object
	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	// Generate the vertex array object
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	//// Generate the Element buffer object
	//var EBO uint32
	//gl.GenBuffers(1, &EBO)
	//gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	return VAO
}

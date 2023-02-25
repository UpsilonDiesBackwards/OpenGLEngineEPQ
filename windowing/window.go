package windowing

import (
	"fmt"
	"github.com/UpsilonDiesBackwards/behngine_epq/input"
	"github.com/UpsilonDiesBackwards/behngine_epq/shaders"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"runtime"
	"strconv"
	"time"
)

// FPS Counter
var startTime = time.Now()
var frameCount int
var FPS float64

var previousTime = time.Now()

var uI = &input.UserInput{} // Create pointer to UserInput struct

var ShouldClose bool
var CursorLocked = true
var cursorLockTimer *time.Timer = nil

func init() {
	runtime.LockOSThread()
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
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// Create a new glfw window called `appWindow`
	appWindow, err := glfw.CreateWindow(
		width,
		height,
		title,
		nil, nil)
	if err != nil {
		fmt.Println("Can not create GLFW window: %s", err)
	}
	appWindow.MakeContextCurrent()
	appWindow.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

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

func EnableFPSCounter(DeltaTime float64) {
	currentTime := time.Now()
	DeltaTime = currentTime.Sub(previousTime).Seconds()
	previousTime = currentTime

	// Handle FPS counter
	frameCount++
	if time.Since(startTime) >= time.Second {
		// Calculate the FPS
		FPS = float64(frameCount) / time.Since(startTime).Seconds()

		// Print the FPS
		str := strconv.FormatFloat(FPS, 'f', 3, 64)
		fmt.Println("FPS: ", str)

		// Reset the frame count and start time
		frameCount = 0
		startTime = time.Now()
	}
}

func EnableWireFrameRendering() {
	// Enable line smoothing and blending
	gl.Enable(gl.BLEND)
	gl.Enable(gl.LINE_SMOOTH)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Hint(gl.LINE_SMOOTH_HINT, gl.NICEST)

	// Set polygon mode to GL_LINE
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
}

func ChangeCursorLockState(appWindow *glfw.Window) {
	if cursorLockTimer == nil || cursorLockTimer.Stop() {
		CursorLocked = !CursorLocked

		if CursorLocked {
			appWindow.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
		} else {
			appWindow.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		}

		cursorLockTimer = time.AfterFunc(500*time.Millisecond, func() {
			cursorLockTimer = nil
		})
	}
}

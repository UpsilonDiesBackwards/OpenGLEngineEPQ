package windowing

import (
	"fmt"
	"github.com/UpsilonDiesBackwards/behngine_epq/input"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"runtime"
	"strconv"
	"time"
)

// FPS Counter
var startTime = time.Now()
var frameCount int
var FPS float64

var previousTime = time.Now()

var userInput = &input.UserInput{} // Create pointer to UserInput struct

var ShouldClose bool
var CursorLocked = true
var cursorLockTimer *time.Timer = nil

func init() {
	runtime.LockOSThread()
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
	// Doesn't properly work atm, spam [ TAB ] for a few seconds then let go, this should toggle
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

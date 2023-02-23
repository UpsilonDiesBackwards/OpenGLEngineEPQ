package windowing

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"strconv"
	"time"
)

// FPS Counter
var startTime = time.Now()
var frameCount int
var FPS float64

var previousTime = time.Now()

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
		str := strconv.FormatFloat(FPS, 'f', 1, 64)
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

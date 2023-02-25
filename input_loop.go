package main

import (
	"fmt"
	"github.com/UpsilonDiesBackwards/behngine_epq/camera"
	"github.com/UpsilonDiesBackwards/behngine_epq/input"
	"github.com/UpsilonDiesBackwards/behngine_epq/windowing"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var ViewportTransform mgl32.Mat4

// ProgramInputLoop During each frame of the program, check for keyboard and mouse input.
func ProgramInputLoop(appWindow *glfw.Window, deltaTime float64, c *camera.Camera, userInput *input.UserInput) error {
	adjustedViewportSpeed := deltaTime * c.Speed

	// Viewport controls
	if input.ActionState[input.VIEWPORT_FORWARD] {
		// Move the camera forward.
		c.Position = c.Position.Add(c.Front.Mul(adjustedViewportSpeed))
	}
	if input.ActionState[input.VIEWPORT_BACKWARDS] {
		// Move the camera backward.
		c.Position = c.Position.Sub(c.Front.Mul(adjustedViewportSpeed))
	}
	if input.ActionState[input.VIEWPORT_LEFT] {
		// Move the camera to the left.
		c.Position = c.Position.Sub(c.Front.Cross(c.Up).Mul(adjustedViewportSpeed))
	}
	if input.ActionState[input.VIEWPORT_RIGHT] {
		// Move the camera to the right.
		c.Position = c.Position.Add(c.Front.Cross(c.Up).Mul(adjustedViewportSpeed))
	}
	if input.ActionState[input.VIEWPORT_RAISE] {
		// Move the camera to the right.
		c.Position = c.Position.Add(c.Up.Mul(adjustedViewportSpeed))
	}
	if input.ActionState[input.VIEWPORT_LOWER] {
		// Move the camera to the right.
		c.Position = c.Position.Sub(c.Up.Mul(adjustedViewportSpeed))
	}

	// Debug and Window controls
	if input.ActionState[input.CHANGE_CURSOR_LOCK_STATE] {
		windowing.ChangeCursorLockState(appWindow)
	}
	if input.ActionState[input.QUIT_PROGRAM] {
		fmt.Println("\nQuitting 3D rendering engine")
		windowing.ShouldClose = true
	}

	// Cursor transform
	ViewportTransform = c.GetTransform()
	uI.CheckpointCursorChange()
	c.UpdateDirection(uI)

	input.CreateInput_Manager()
	return nil
}

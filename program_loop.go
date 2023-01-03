package main

import (
	"fmt"
	"github.com/UpsilonDiesBackwards/behngine_epq/input"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var iUserInput = input.UserInput{}

func ProgLoop(appWindow *glfw.Window) error {
	if input.ActionState[input.INPUT_TEST] {
		fmt.Println("\nInput test!")
	}

	if input.ActionState[input.QUIT_PROGRAM] {
		fmt.Println("\nQuitting 3D rendering engine")
		appWindow.SetShouldClose(true)
	}

	fmt.Println("Cursor: ", iUserInput.CursorChange())

	return nil
}

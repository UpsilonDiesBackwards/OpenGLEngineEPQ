package main

import (
	"fmt"
	"github.com/UpsilonDiesBackwards/behngine_epq/input"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func ProgLoop(appWindow *glfw.Window) error {
	if input.ActionState[input.INPUT_TEST] {
		fmt.Println("\nInput test!")
	}

	if input.ActionState[input.QUIT_PROGRAM] {
		fmt.Println("\nQuitting 3D rendering engine")
		appWindow.SetShouldClose(true)
	}

	input.Input_Manager(appWindow)
	appWindow.SetCursorPosCallback(input.UserInput{}.MouseCallBack)
	//fmt.Println("Mouse: ", input.UserInput{})

	return nil
}

package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

func ProcessInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}

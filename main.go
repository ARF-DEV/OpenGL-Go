package main

import (
	"fmt"
	"learn-open-gl/callbacks"
	"learn-open-gl/gogl"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		fmt.Println("Error initializing glfw: ", err.Error())
		return
	}

	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(800, 600, "Test", nil, nil)
	if err != nil {
		fmt.Println("Error creating window: ", err.Error())
		return
	}

	defer window.Destroy()

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		fmt.Println("Error initializing gl: ", err.Error())
		return
	}
	gl.Viewport(0, 0, 800, 600)

	window.SetFramebufferSizeCallback(callbacks.FrameBufferSizeCallback)
	window.SetKeyCallback(callbacks.KeyCallback)

	// vertexShader, err := gogl.LoadShader("./shaders/vertex.vert", gl.VERTEX_SHADER)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// fragmentShader, err := gogl.LoadShader("./shaders/fragment.frag", gl.FRAGMENT_SHADER)
	// if err != nil {
	// 	panic(err.Error())
	// }

	shaderProgram, err := gogl.CreateProgramStructFromPaths(
		[]uint32{gl.VERTEX_SHADER, gl.FRAGMENT_SHADER},
		"./shaders/vertex.vert", "./shaders/fragment.frag",
	)
	if err != nil {
		panic(err.Error())
	}

	vertices := []float32{
		0.5, 0.5, 0.0, // top right
		0.5, -0.5, 0.0, // bottom right
		-0.5, -0.5, 0.0, // bottom let
		-0.5, 0.5, 0.0, // top left

	}

	indices := []uint32{ // note that we start from 0!
		0, 1, 3, // first triangle
		1, 2, 3} // second triangle

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	var EBO uint32
	gl.GenBuffers(1, &EBO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.UseProgram(shaderProgram.ID)
	for !window.ShouldClose() {
		// input

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.BindVertexArray(VAO)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()
		if shaderProgram.IsUpdated() {
			shaderProgram.ReloadProgram()
			gl.UseProgram(shaderProgram.ID)
		}
	}

}

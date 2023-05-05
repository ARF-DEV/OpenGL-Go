package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"learn-open-gl/callbacks"
	"learn-open-gl/gogl"
	"log"
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

	shaderProgram, err := gogl.CreateShader(
		"./shaders/vertex.vs", "./shaders/fragment.fs",
	)
	if err != nil {
		panic(err.Error())
	}

	vertices := []float32{
		// positions          // colors           // texture coords
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // top right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom left
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // top left

	}

	indices := []uint32{ // note that we start from 0!
		0, 1, 3, // first triangle
		1, 2, 3} // second triangle

	boxTexture, err := gogl.CreateTextureFromFile("images/container.jpg")
	if err != nil {
		log.Println("Failed to create texture: ", err.Error())
		return
	}

	faceTexture, err := gogl.CreateTextureFromFile("images/awesomeface.png")
	if err != nil {
		log.Println("Failed to create texture: ", err.Error())
		return
	}

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

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 8*4, 3*4)
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, 8*4, 6*4)
	gl.EnableVertexAttribArray(2)

	gl.UseProgram(shaderProgram.ID)
	gl.Uniform1i(gl.GetUniformLocation(shaderProgram.ID, gl.Str("texture1\x00")), 0)
	gl.Uniform1i(gl.GetUniformLocation(shaderProgram.ID, gl.Str("texture2\x00")), 1)

	for !window.ShouldClose() {
		shaderProgram.ReloadOnUpdate()
		// input

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, boxTexture.ID)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, faceTexture.ID)

		gl.BindVertexArray(VAO)
		gl.UseProgram(shaderProgram.ID)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()
		// if shaderProgram.IsUpdated() {
		// 	fmt.Println("BERUBAH")
		// 	shaderProgram.ReloadProgram()
		// }
		// if shaderProgram.IsUpdated() {
		// 	shaderProgram
		// }
		err := gl.GetError()
		if err != gl.NO_ERROR {
			log.Println(err)
			panic(err)
		}
	}

}

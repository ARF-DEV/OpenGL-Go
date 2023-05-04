package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"learn-open-gl/callbacks"
	"learn-open-gl/gogl"
	"log"
	"os"
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

	borderColor := []float32{
		1.0, 1.0, 0.0, 1.0,
	}

	imageFile, _ := os.Open("./images/container.jpg")
	img, _, err := image.Decode(imageFile)
	if err != nil {
		panic("Failed to decode image: " + err.Error())
	}
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Rect, img, image.Pt(0, 0), draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.MIRRORED_REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.MIRRORED_REPEAT)
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, (*float32)(gl.Ptr(borderColor)))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Bounds().Size().X), int32(img.Bounds().Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	imageFileFace, _ := os.Open("./images/awesomeface.png")
	img2, _, err := image.Decode(imageFileFace)
	if err != nil {
		panic("Failed to decode image: " + err.Error())
	}
	rgba2 := image.NewRGBA(img2.Bounds())
	draw.Draw(rgba2, rgba2.Rect, img2, image.Pt(0, 0), draw.Src)

	var texture2 uint32
	gl.GenTextures(1, &texture2)
	gl.BindTexture(gl.TEXTURE_2D, texture2)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.MIRRORED_REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.MIRRORED_REPEAT)
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, (*float32)(gl.Ptr(borderColor)))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Bounds().Size().X), int32(img.Bounds().Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba2.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

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

		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)

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

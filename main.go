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
	"github.com/go-gl/mathgl/mgl32"
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
		// positions  // texture coords
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
	}

	indices := []uint32{ // note that we start from 0!
		0, 1, 3, // first triangle
		1, 2, 3} // second triangle

	boxTexture, err := gogl.CreateTextureFromFile("images/container.jpg", gl.REPEAT, gl.LINEAR, gl.LINEAR_MIPMAP_LINEAR, false)
	if err != nil {
		log.Println("Failed to create texture: ", err.Error())
		return
	}

	faceTexture, err := gogl.CreateTextureFromFile("images/awesomeface.png", gl.REPEAT, gl.NEAREST, gl.NEAREST, true)
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

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)

	gl.UseProgram(shaderProgram.ID)
	gl.Uniform1i(gl.GetUniformLocation(shaderProgram.ID, gl.Str("texture1\x00")), 0)
	gl.Uniform1i(gl.GetUniformLocation(shaderProgram.ID, gl.Str("texture2\x00")), 1)

	// Vector stuff
	// {
	// 	vec := mgl32.Vec4{1.0, 0.0, 0.0, 1.0}
	// 	trans := mgl32.Translate3D(1.0, 1.0, 0.0)
	// 	result := trans.Mul4x1(vec)
	// 	fmt.Printf("%f, %f, %f\n", result.X(), result.Y(), result.Z())
	// }

	// Model matrix is where the position, scale, and rotation of the object is defined
	view := mgl32.Translate3D(0, 0, -3)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), 800.0/600.0, 0.1, 100.0)
	projectionLoc := gl.GetUniformLocation(shaderProgram.ID, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionLoc, 1, false, &projection[0])

	cubePositions := []mgl32.Vec3{
		{0.0, 0.0, 0.0},
		{2.0, 5.0, -15.0},
		{-1.5, -2.2, -2.5},
		{-3.8, -2.0, -12.3},
		{2.4, -0.4, -3.5},
		{-1.7, 3.0, -7.5},
		{1.3, -2.0, -2.5},
		{1.5, 2.0, -2.5},
		{1.5, 0.2, -1.5},
		{-1.3, 1.0, -1.5},
	}

	gl.Enable(gl.DEPTH_TEST)
	for !window.ShouldClose() {
		shaderProgram.ReloadOnUpdate()
		// input

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update Position

		modelLoc := gl.GetUniformLocation(shaderProgram.ID, gl.Str("model\x00"))
		viewLoc := gl.GetUniformLocation(shaderProgram.ID, gl.Str("view\x00"))

		gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])

		gl.ActiveTexture(gl.TEXTURE0)
		boxTexture.Bind()
		gl.ActiveTexture(gl.TEXTURE1)
		faceTexture.Bind()

		gl.BindVertexArray(VAO)
		gl.UseProgram(shaderProgram.ID)
		// gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
		for i, pos := range cubePositions {
			posMat := mgl32.Translate3D(pos.X(), pos.Y(), pos.Z())
			var rotMat mgl32.Mat4
			if i%3 == 0 {
				rotMat = mgl32.HomogRotate3D(mgl32.DegToRad(float32(i)*20.0)+float32(glfw.GetTime()), mgl32.Vec3{1.0, 0.3, 0.5}.Normalize())
			} else {
				rotMat = mgl32.HomogRotate3D(mgl32.DegToRad(float32(i)*20.0), mgl32.Vec3{1.0, 0.3, 0.5}.Normalize())
			}
			model := posMat.Mul4(rotMat)
			gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()

		err := gl.GetError()
		if err != gl.NO_ERROR {
			log.Println(err)
			panic(err)
		}
	}

}

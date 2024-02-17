// package main

// import (
// 	"fmt"
// 	_ "image/jpeg"
// 	_ "image/png"
// 	"learn-open-gl/callbacks"
// 	"learn-open-gl/engine"
// 	"learn-open-gl/gogl"
// 	"log"
// 	"runtime"

// 	"github.com/go-gl/gl/v3.3-core/gl"
// 	"github.com/go-gl/glfw/v3.3/glfw"
// 	"github.com/go-gl/mathgl/mgl32"
// )

// var deltaTime = 0.0
// var lastFrame = 0.0

// func main() {
// 	runtime.LockOSThread()

// 	if err := glfw.Init(); err != nil {
// 		fmt.Println("Error initializing glfw: ", err.Error())
// 		return
// 	}

// 	defer glfw.Terminate()

// 	glfw.WindowHint(glfw.ContextVersionMajor, 3)
// 	glfw.WindowHint(glfw.ContextVersionMinor, 3)
// 	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

// 	window, err := glfw.CreateWindow(800, 600, "Test", nil, nil)
// 	if err != nil {
// 		fmt.Println("Error creating window: ", err.Error())
// 		return
// 	}

// 	defer window.Destroy()

// 	window.MakeContextCurrent()

// 	if err := gl.Init(); err != nil {
// 		fmt.Println("Error initializing gl: ", err.Error())
// 		return
// 	}
// 	gl.Viewport(0, 0, 800, 600)

// 	window.SetFramebufferSizeCallback(callbacks.FrameBufferSizeCallback)
// 	window.SetKeyCallback(callbacks.KeyCallback)
// 	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

// 	// Camera thing
// 	fov := 45.0
// 	mainCamera := engine.CreateDefaultCamera()
// 	firstMouse := true
// 	var xLast float64 = 300
// 	var yLast float64 = 400

// 	window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {

// 		if firstMouse {
// 			xLast = xpos
// 			yLast = ypos
// 			firstMouse = false
// 		}
// 		xOffset := xpos - xLast
// 		yOffset := yLast - ypos
// 		xLast = xpos
// 		yLast = ypos

// 		mainCamera.ProcessMouseMovement(float32(xOffset), float32(yOffset))
// 	})

// 	window.SetScrollCallback(func(w *glfw.Window, xoff, yoff float64) {
// 		fov -= yoff
// 		if fov < 1 {
// 			fov = 1
// 		}
// 		if fov > 45 {
// 			fov = 45
// 		}
// 	})

// 	shaderProgram, err := gogl.CreateShader(
// 		"./shaders/vertex.vs", "./shaders/fragment.fs",
// 	)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	lightShaderProgram, err := gogl.CreateShader(
// 		"./shaders/vertex1.vs", "./shaders/fragment1.fs",
// 	)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	vertices := []float32{
// 		// positions  // texture coords
// 		-0.5, -0.5, -0.5, 0.0, 0.0,
// 		0.5, -0.5, -0.5, 1.0, 0.0,
// 		0.5, 0.5, -0.5, 1.0, 1.0,
// 		0.5, 0.5, -0.5, 1.0, 1.0,
// 		-0.5, 0.5, -0.5, 0.0, 1.0,
// 		-0.5, -0.5, -0.5, 0.0, 0.0,

// 		-0.5, -0.5, 0.5, 0.0, 0.0,
// 		0.5, -0.5, 0.5, 1.0, 0.0,
// 		0.5, 0.5, 0.5, 1.0, 1.0,
// 		0.5, 0.5, 0.5, 1.0, 1.0,
// 		-0.5, 0.5, 0.5, 0.0, 1.0,
// 		-0.5, -0.5, 0.5, 0.0, 0.0,

// 		-0.5, 0.5, 0.5, 1.0, 0.0,
// 		-0.5, 0.5, -0.5, 1.0, 1.0,
// 		-0.5, -0.5, -0.5, 0.0, 1.0,
// 		-0.5, -0.5, -0.5, 0.0, 1.0,
// 		-0.5, -0.5, 0.5, 0.0, 0.0,
// 		-0.5, 0.5, 0.5, 1.0, 0.0,

// 		0.5, 0.5, 0.5, 1.0, 0.0,
// 		0.5, 0.5, -0.5, 1.0, 1.0,
// 		0.5, -0.5, -0.5, 0.0, 1.0,
// 		0.5, -0.5, -0.5, 0.0, 1.0,
// 		0.5, -0.5, 0.5, 0.0, 0.0,
// 		0.5, 0.5, 0.5, 1.0, 0.0,

// 		-0.5, -0.5, -0.5, 0.0, 1.0,
// 		0.5, -0.5, -0.5, 1.0, 1.0,
// 		0.5, -0.5, 0.5, 1.0, 0.0,
// 		0.5, -0.5, 0.5, 1.0, 0.0,
// 		-0.5, -0.5, 0.5, 0.0, 0.0,
// 		-0.5, -0.5, -0.5, 0.0, 1.0,

// 		-0.5, 0.5, -0.5, 0.0, 1.0,
// 		0.5, 0.5, -0.5, 1.0, 1.0,
// 		0.5, 0.5, 0.5, 1.0, 0.0,
// 		0.5, 0.5, 0.5, 1.0, 0.0,
// 		-0.5, 0.5, 0.5, 0.0, 0.0,
// 		-0.5, 0.5, -0.5, 0.0, 1.0,
// 	}

// 	indices := []uint32{ // note that we start from 0!
// 		0, 1, 3, // first triangle
// 		1, 2, 3} // second triangle

// 	// boxTexture, err := gogl.CreateTextureFromFile("images/container.jpg", gl.REPEAT, gl.LINEAR, gl.LINEAR_MIPMAP_LINEAR, false)
// 	// if err != nil {
// 	// 	log.Println("Failed to create texture: ", err.Error())
// 	// 	return
// 	// }

// 	// faceTexture, err := gogl.CreateTextureFromFile("images/awesomeface.png", gl.REPEAT, gl.NEAREST, gl.NEAREST, true)
// 	// if err != nil {
// 	// 	log.Println("Failed to create texture: ", err.Error())
// 	// 	return
// 	// }

// 	var VAO uint32
// 	gl.GenVertexArrays(1, &VAO)
// 	gl.BindVertexArray(VAO)

// 	var EBO uint32
// 	gl.GenBuffers(1, &EBO)
// 	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
// 	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

// 	var VBO uint32
// 	gl.GenBuffers(1, &VBO)
// 	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
// 	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

// 	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
// 	gl.EnableVertexAttribArray(0)

// 	// gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
// 	// gl.EnableVertexAttribArray(1)

// 	// gl.UseProgram(shaderProgram.ID)
// 	// gl.Uniform1i(gl.GetUniformLocation(shaderProgram.ID, gl.Str("texture1\x00")), 0)
// 	// gl.Uniform1i(gl.GetUniformLocation(shaderProgram.ID, gl.Str("texture2\x00")), 1)
// 	// gl.BindVertexArray(0)

// 	var lightVBO uint32
// 	gl.GenBuffers(1, &lightVBO)
// 	gl.BindBuffer(gl.ARRAY_BUFFER, lightVBO)
// 	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

// 	var lightVAO uint32
// 	gl.GenVertexArrays(1, &lightVAO)
// 	gl.BindVertexArray(lightVAO)
// 	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 5*4, 0)
// 	// gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
// 	gl.EnableVertexAttribArray(0)
// 	gl.EnableVertexAttribArray(1)

// 	gl.BindVertexArray(0)

// 	view := mgl32.Translate3D(0, 0, -3)

// 	projection := mgl32.Perspective(mgl32.DegToRad(float32(fov)), 800.0/600.0, 0.1, 100.0)
// 	projectionLoc := gl.GetUniformLocation(shaderProgram.ID, gl.Str("projection\x00"))
// 	gl.UniformMatrix4fv(projectionLoc, 1, false, &projection[0])

// 	gl.UseProgram(shaderProgram.ID)
// 	objectPosition := mgl32.Translate3D(0, 0, 0)
// 	modelLoc := gl.GetUniformLocation(shaderProgram.ID, gl.Str("model\x00"))
// 	gl.UniformMatrix4fv(modelLoc, 1, false, &objectPosition[0])
// 	colorLoc := gl.GetUniformLocation(shaderProgram.ID, gl.Str("objectColor\x00"))
// 	gl.Uniform3f(colorLoc, 1.0, 0.5, 0.31)

// 	lightColorLoc := gl.GetUniformLocation(shaderProgram.ID, gl.Str("lightColor\x00"))
// 	gl.Uniform3f(lightColorLoc, 1.0, 1.0, 1.0)
// 	gl.UseProgram(0)

// 	// set for light shader
// 	gl.UseProgram(lightShaderProgram.ID)

// 	lightPosition := mgl32.Translate3D(0, 0, 3)
// 	modelLoc = gl.GetUniformLocation(lightShaderProgram.ID, gl.Str("model\x00"))
// 	gl.UniformMatrix4fv(modelLoc, 1, false, &lightPosition[0])

// 	projection = mgl32.Perspective(mgl32.DegToRad(float32(fov)), 800.0/600.0, 0.1, 100.0)
// 	projectionLoc = gl.GetUniformLocation(lightShaderProgram.ID, gl.Str("projection\x00"))
// 	gl.UniformMatrix4fv(projectionLoc, 1, false, &projection[0])

// 	// fmt.Println(gl.GetError())
// 	// Creating Camera View Matrix
// 	// cameraPos := mgl32.Vec3{0, 0, 3}
// 	// cameraTarget := mgl32.Vec3{0, 0, 0}
// 	// The name direction vector is not the best chosen name, since it is actually pointing in the reverse direction of what it is targeting.
// 	// https://stackoverflow.com/questions/60362629/confusion-on-opengls-camera-and-camera-space-transformation
// 	// cameraDirection := cameraPos.Sub(cameraTarget).Normalize()
// 	// cameraRight := mgl32.Vec3{0, 1, 0}.Cross(cameraDirection)
// 	// cameraUp := cameraDirection.Cross(cameraRight)

// 	//

// 	gl.Enable(gl.DEPTH_TEST)

// 	for !window.ShouldClose() {
// 		currentFrame := glfw.GetTime()
// 		deltaTime = currentFrame - lastFrame
// 		lastFrame = currentFrame

// 		shaderProgram.ReloadOnUpdate()
// 		// input
// 		if window.GetKey(glfw.KeyW) == glfw.Press {
// 			mainCamera.ProcessKeyboard(engine.FORWARD, float32(deltaTime))
// 		}
// 		if window.GetKey(glfw.KeyS) == glfw.Press {
// 			mainCamera.ProcessKeyboard(engine.BACKWARD, float32(deltaTime))
// 		}
// 		if window.GetKey(glfw.KeyA) == glfw.Press {
// 			mainCamera.ProcessKeyboard(engine.LEFT, float32(deltaTime))
// 		}
// 		if window.GetKey(glfw.KeyD) == glfw.Press {
// 			mainCamera.ProcessKeyboard(engine.RIGHT, float32(deltaTime))
// 		}

// 		if window.GetKey(glfw.KeyEscape) == glfw.Press {
// 			window.SetShouldClose(true)
// 		}

// 		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
// 		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

// 		// Update Position

// 		viewLoc := gl.GetUniformLocation(shaderProgram.ID, gl.Str("view\x00"))
// 		view = mainCamera.GetLookUpMatrix()
// 		gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])

// 		// gl.ActiveTexture(gl.TEXTURE0)
// 		// boxTexture.Bind()
// 		// gl.ActiveTexture(gl.TEXTURE1)
// 		// faceTexture.Bind()

// 		gl.BindVertexArray(VAO)
// 		gl.UseProgram(shaderProgram.ID)
// 		// gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
// 		gl.DrawArrays(gl.TRIANGLES, 0, 36)

// 		// set for light shader
// 		viewLoc = gl.GetUniformLocation(lightShaderProgram.ID, gl.Str("view\x00"))
// 		view = mainCamera.GetLookUpMatrix()
// 		gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])

// 		gl.BindVertexArray(lightVAO)
// 		gl.UseProgram(lightShaderProgram.ID)
// 		gl.DrawArrays(gl.TRIANGLES, 0, 36)

// 		gl.BindVertexArray(0)

// 		window.SwapBuffers()
// 		glfw.PollEvents()

// 		err := gl.GetError()
// 		if err != gl.NO_ERROR {
// 			log.Println(err)
// 			panic(err)
// 		}
// 	}
// }

package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"learn-open-gl/callbacks"
	"learn-open-gl/engine"
	"learn-open-gl/gogl"
	"log"
	"math"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var deltaTime = 0.0
var lastFrame = 0.0

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
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// Camera thing
	fov := 45.0
	mainCamera := engine.CreateDefaultCamera()
	firstMouse := true
	var xLast float64 = 300
	var yLast float64 = 400

	window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {

		if firstMouse {
			xLast = xpos
			yLast = ypos
			firstMouse = false
		}
		xOffset := xpos - xLast
		yOffset := yLast - ypos
		xLast = xpos
		yLast = ypos

		mainCamera.ProcessMouseMovement(float32(xOffset), float32(yOffset))
	})

	window.SetScrollCallback(func(w *glfw.Window, xo, yoff float64) {
		fov -= yoff
		if fov < 1 {
			fov = 1
		}
		if fov > 45 {
			fov = 45
		}
	})

	shaderProgram, err := gogl.CreateShader(
		"./shaders/vertex.vs", "./shaders/fragmentnew.fs",
	)
	if err != nil {
		panic(err.Error())
	}
	lightShaderProgram, err := gogl.CreateShader(
		"./shaders/vertex.vs", "./shaders/fragment1.fs",
	)
	if err != nil {
		panic(err.Error())
	}

	vertices := []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,

		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,
		-0.5, 0.5, -0.5, -1.0, 0.0, 0.0, 1.0, 1.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, -1.0, 0.0, 0.0, 0.0, 0.0,
		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 1.0, 1.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
	}

	indices := []uint32{ // note that we start from 0!
		0, 1, 3, // first triangle
		1, 2, 3} // second triangle

	cubePosition := []mgl32.Vec3{
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

	diffuseMap, err := gogl.CreateTextureFromFile("images/container2.png", gl.CLAMP_TO_BORDER, gl.NEAREST, gl.NEAREST, true)
	if err != nil {
		panic(err)
	}
	specularMap, err := gogl.CreateTextureFromFile("images/container2_specular.png", gl.CLAMP_TO_BORDER, gl.NEAREST, gl.NEAREST, true)
	if err != nil {
		panic(err)
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

	var lightVBO uint32
	gl.GenBuffers(1, &lightVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, lightVBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	var lightVAO uint32
	gl.GenVertexArrays(1, &lightVAO)
	gl.BindVertexArray(lightVAO)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 8*4, 0)
	// gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(0)
	// gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	view := mgl32.Translate3D(0, 0, -3)

	gl.UseProgram(shaderProgram.ID)
	projection := mgl32.Perspective(mgl32.DegToRad(float32(fov)), 800.0/600.0, 0.1, 100.0)
	projectionLoc := gl.GetUniformLocation(shaderProgram.ID, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionLoc, 1, false, &projection[0])
	checkError("OAPAA")

	objectPosition := mgl32.Translate3D(0, 0, 0)
	modelLoc := gl.GetUniformLocation(shaderProgram.ID, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelLoc, 1, false, &objectPosition[0])
	colorLoc := gl.GetUniformLocation(shaderProgram.ID, gl.Str("objectColor\x00"))
	gl.Uniform3f(colorLoc, 1.0, 0.5, 0.31)
	checkError("aowkdoakwdo")

	lightPos := mgl32.Vec3{2, 0, 2.0}
	lightTransform := mgl32.Translate3D(lightPos[0], lightPos[1], lightPos[2]).Mul4(mgl32.Scale3D(0.2, 0.2, 0.2))

	shaderProgram.SetInt("material.diffuse", 0)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, diffuseMap.ID)
	shaderProgram.SetInt("material.specular", 1)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, specularMap.ID)
	shaderProgram.SetUniformFloat("material.shininess", 64)
	checkError("awdklawpodkl[]")

	// store light position
	lightColor := mgl32.Vec3{1, 1, 1}
	diffuseColor := lightColor.Mul(0.5)
	ambientColor := lightColor.Mul(0.2)
	shaderProgram.SetVec3("dirlight.ambient", ambientColor)
	shaderProgram.SetVec3("dirlight.specular", mgl32.Vec3{0.5, 0.5, 0.5})
	shaderProgram.SetVec3("dirlight.diffuse", diffuseColor)
	shaderProgram.SetVec3("dirlight.direction", mgl32.Vec3{-0.2, -1, -0.3})
	lightPositions := []mgl32.Vec3{
		{0.7, 0.2, 2.0},
		{2.3, -3.3, -4.0},
		{-4.0, 2.0, -12.0},
		{0.0, 0.0, -3.0},
	}
	for i, pos := range lightPositions {
		shaderProgram.SetVec3(fmt.Sprintf("pointlights[%d].position", i), pos)
		shaderProgram.SetVec3(fmt.Sprintf("pointlights[%d].ambient", i), mgl32.Vec3{0.5, 0.5, 0.5})
		shaderProgram.SetVec3(fmt.Sprintf("pointlights[%d].diffuse", i), mgl32.Vec3{0.5, 0.5, 0.5})
		shaderProgram.SetVec3(fmt.Sprintf("pointlights[%d].specular", i), mgl32.Vec3{1, 1, 1})
		shaderProgram.SetUniformFloat(fmt.Sprintf("pointlights[%d].constant", i), 1)
		shaderProgram.SetUniformFloat(fmt.Sprintf("pointlights[%d].linear", i), 0.09)
		shaderProgram.SetUniformFloat(fmt.Sprintf("pointlights[%d].quadratic", i), 0.032)
	}

	gl.UseProgram(0)
	// set for light shader
	gl.UseProgram(lightShaderProgram.ID)
	modelLoc = gl.GetUniformLocation(lightShaderProgram.ID, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelLoc, 1, false, &lightTransform[0])

	projection = mgl32.Perspective(mgl32.DegToRad(float32(fov)), 800.0/600.0, 0.1, 100.0)
	projectionLoc = gl.GetUniformLocation(lightShaderProgram.ID, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionLoc, 1, false, &projection[0])

	checkError("HAUUIU")
	// fmt.Println(gl.GetError())
	// Creating Camera View Matrix
	// cameraPos := mgl32.Vec3{0, 0, 3}
	// cameraTarget := mgl32.Vec3{0, 0, 0}
	// The name direction vector is not the best chosen name, since it is actually pointing in the reverse direction of what it is targeting.
	// https://stackoverflow.com/questions/60362629/confusion-on-opengls-camera-and-camera-space-transformation
	// cameraDirection := cameraPos.Sub(cameraTarget).Normalize()
	// cameraRight := mgl32.Vec3{0, 1, 0}.Cross(cameraDirection)
	// cameraUp := cameraDirection.Cross(cameraRight)

	//

	gl.Enable(gl.DEPTH_TEST)

	for !window.ShouldClose() {
		currentFrame := glfw.GetTime()
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		shaderProgram.ReloadOnUpdate()
		// input
		if window.GetKey(glfw.KeyW) == glfw.Press {
			mainCamera.ProcessKeyboard(engine.FORWARD, float32(deltaTime))
		}
		if window.GetKey(glfw.KeyS) == glfw.Press {
			mainCamera.ProcessKeyboard(engine.BACKWARD, float32(deltaTime))
		}
		if window.GetKey(glfw.KeyA) == glfw.Press {
			mainCamera.ProcessKeyboard(engine.LEFT, float32(deltaTime))
		}
		if window.GetKey(glfw.KeyD) == glfw.Press {
			mainCamera.ProcessKeyboard(engine.RIGHT, float32(deltaTime))
		}
		if window.GetKey(glfw.KeySpace) == glfw.Press {
			mainCamera.ProcessKeyboard(engine.UP, float32(deltaTime))
		}
		if window.GetKey(glfw.KeyLeftControl) == glfw.Press {
			mainCamera.ProcessKeyboard(engine.DOWN, float32(deltaTime))
		}

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update Position
		gl.UseProgram(shaderProgram.ID)
		cameraPosLoc := gl.GetUniformLocation(shaderProgram.ID, gl.Str("cameraPosition\x00"))
		gl.Uniform3fv(cameraPosLoc, 1, &mainCamera.Position[0])

		viewLoc := gl.GetUniformLocation(shaderProgram.ID, gl.Str("view\x00"))
		view = mainCamera.GetLookUpMatrix()
		gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])

		// lightPos = mgl32.Vec3{2 * float32(math.Sin(glfw.GetTime())), 0, 2 * float32(math.Cos(glfw.GetTime()))}
		// lightPosLoc := gl.GetUniformLocation(shaderProgram.ID, gl.Str("lightPos\x00"))
		// gl.Uniform3fv(lightPosLoc, 1, &lightPos[0])
		// fmt.Println(lightPos)

		// lightColor[0] = float32(math.Sin(glfw.GetTime() * 2))
		// lightColor[1] = float32(math.Sin(glfw.GetTime() * 0.7))
		// lightColor[2] = float32(math.Sin(glfw.GetTime() * 1.3))
		// diffuseColor := lightColor.Mul(0.5)
		// ambientColor := lightColor.Mul(0.2)
		// shaderProgram.SetVec3("light.ambient", ambientColor)
		// shaderProgram.SetVec3("light.specular", mgl32.Vec3{0.5, 0.5, 0.5})
		// shaderProgram.SetVec3("light.diffuse", diffuseColor)

		// gl.ActiveTexture(gl.TEXTURE0)
		// boxTexture.Bind()
		// gl.ActiveTexture(gl.TEXTURE1)
		// faceTexture.Bind()
		shaderProgram.SetVec3("spotlight.position", mainCamera.Position)
		shaderProgram.SetVec3("spotlight.direction", mainCamera.Front)
		shaderProgram.SetVec3("spotlight.ambient", ambientColor)
		shaderProgram.SetVec3("spotlight.specular", mgl32.Vec3{0.5, 0.5, 0.5})
		shaderProgram.SetVec3("spotlight.diffuse", diffuseColor)
		shaderProgram.SetUniformFloat("spotlight.constant", 1)
		shaderProgram.SetUniformFloat("spotlight.linear", 0.09)
		shaderProgram.SetUniformFloat("spotlight.quadratic", 0.032)
		shaderProgram.SetUniformFloat("spotlight.innerCutoff", float32(math.Cos(float64(mgl32.DegToRad(12.5)))))
		shaderProgram.SetUniformFloat("spotlight.outerCutoff", float32(math.Cos(float64(mgl32.DegToRad(17.5)))))

		gl.BindVertexArray(VAO)
		// gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
		for index, pos := range cubePosition {
			angle := index * 20.0
			translate := mgl32.Translate3D(pos[0], pos[1], pos[2])
			rotate := mgl32.HomogRotate3D(mgl32.DegToRad(float32(angle)), mgl32.Vec3{1, 0.3, 0.5}.Normalize())

			model := translate.Mul4(rotate)
			modelLoc := gl.GetUniformLocation(shaderProgram.ID, gl.Str("model\x00"))
			gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		gl.UseProgram(lightShaderProgram.ID)
		// set for light shader
		viewLoc = gl.GetUniformLocation(lightShaderProgram.ID, gl.Str("view\x00"))
		view = mainCamera.GetLookUpMatrix()
		gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])

		for _, pos := range lightPositions {
			lightTransform = mgl32.Translate3D(pos[0], pos[1], pos[2]).Mul4(mgl32.Scale3D(0.2, 0.2, 0.2))
			fmt.Println(pos)
			modelLoc = gl.GetUniformLocation(lightShaderProgram.ID, gl.Str("model\x00"))
			gl.UniformMatrix4fv(modelLoc, 1, false, &lightTransform[0])
			gl.BindVertexArray(lightVAO)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()

		checkError("UHUY")
	}
}
func checkError(tag string) {
	err := gl.GetError()
	if err != gl.NO_ERROR {
		fmt.Println(tag)
		log.Println(err)
		panic(err)
	}
}

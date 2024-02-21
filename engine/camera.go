package engine

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

const YAW = -90.0
const PITCH = 0.0
const SPEED = 10
const SENSITIVITY = 0.1
const ZOOM = 45.0

type CameraMovement int32

const (
	FORWARD CameraMovement = iota
	BACKWARD
	LEFT
	RIGHT
	UP
	DOWN
)

type Camera struct {
	Position mgl32.Vec3
	Front    mgl32.Vec3
	Up       mgl32.Vec3
	Right    mgl32.Vec3
	WorldUp  mgl32.Vec3

	Yaw   float32
	Pitch float32

	MovementSpeed    float32
	MouseSensitivity float32
	Zoom             float32
}

func CreateNewCamera(Position, Front, Up, WorldUp mgl32.Vec3, Yaw, Pitch, MovementSpeed, MouseSensitivity, Zoom float32) Camera {
	return Camera{
		Position:         Position,
		Front:            Front,
		WorldUp:          WorldUp,
		Yaw:              Yaw,
		Pitch:            Pitch,
		MovementSpeed:    MovementSpeed,
		MouseSensitivity: MouseSensitivity,
		Zoom:             Zoom,
	}
}

func CreateDefaultCamera() Camera {
	return Camera{
		Position:         mgl32.Vec3{0, 0, 3},
		Front:            mgl32.Vec3{0, 0, -1},
		WorldUp:          mgl32.Vec3{0, 1, 0},
		Yaw:              YAW,
		Pitch:            PITCH,
		MovementSpeed:    SPEED,
		MouseSensitivity: SENSITIVITY,
		Zoom:             ZOOM,
	}
}

func (c *Camera) GetLookUpMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Position.Add(c.Front), c.WorldUp)
}

func (c *Camera) ProcessKeyboard(direction CameraMovement, deltaTime float32) {
	velocity := c.MovementSpeed * deltaTime
	if direction == FORWARD {
		c.Position = c.Position.Add(c.Front.Mul(velocity))
	}
	if direction == BACKWARD {
		c.Position = c.Position.Sub(c.Front.Mul(velocity))
	}
	if direction == RIGHT {
		c.Position = c.Position.Add(c.Front.Cross(c.WorldUp).Normalize().Mul(velocity))
	}
	if direction == LEFT {
		c.Position = c.Position.Sub(c.Front.Cross(c.WorldUp).Normalize().Mul(velocity))
	}
	if direction == UP {
		c.Position = c.Position.Add(c.WorldUp.Mul(velocity))
	}
	if direction == DOWN {
		c.Position = c.Position.Sub(c.WorldUp.Mul(velocity))
	}
}

func (c *Camera) ProcessMouseMovement(xOffset, yOffset float32) {
	xOffset *= c.MouseSensitivity
	yOffset *= c.MouseSensitivity

	c.Yaw += xOffset
	c.Pitch += yOffset

	if c.Pitch > 89 {
		c.Pitch = 89
	}
	if c.Pitch < -89 {
		c.Pitch = -89
	}

	c.updateVector()
}

func (c *Camera) updateVector() {
	var direction mgl32.Vec3 = mgl32.Vec3{}
	radPitch := mgl32.DegToRad(float32(c.Yaw))
	radYaw := mgl32.DegToRad(float32(c.Pitch))

	direction[0] = float32(math.Cos(float64(radPitch)) * math.Cos(float64(radYaw)))
	direction[1] = float32(math.Sin(float64(radYaw)))
	direction[2] = float32(math.Sin(float64(radPitch)) * math.Cos(float64(radYaw)))

	c.Front = direction.Normalize()
	c.Right = c.Front.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Front).Normalize()
}

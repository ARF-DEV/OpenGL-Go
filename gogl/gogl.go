package gogl

import (
	"errors"
	"os"

	"github.com/go-gl/gl/v3.3-core/gl"
)

func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}

func CreateShader(shaderSource string, shaderType uint32) (uint32, error) {
	shaderId := gl.CreateShader(shaderType)
	shaderSource += "\x00"
	cSource, free := gl.Strs(shaderSource)
	gl.ShaderSource(shaderId, 1, cSource, nil)
	free()
	gl.CompileShader(shaderId)
	var glErr int32
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &glErr)
	if glErr == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLength)
		log := string(make([]byte, logLength+1))
		gl.GetShaderInfoLog(shaderId, logLength, nil, gl.Str(log))
		return 0, errors.New("Error Failed to compile shader:\n" + log)
	}
	return shaderId, nil
}

func CreateProgram(shadersId ...uint32) (uint32, error) {
	programId := gl.CreateProgram()
	if err := LinkProgram(programId, shadersId...); err != nil {
		return 0, err
	}
	return programId, nil

}

func LinkProgram(programId uint32, shadersId ...uint32) error {
	for _, shader := range shadersId {
		gl.AttachShader(programId, shader)
	}
	gl.LinkProgram(programId)

	var glErr int32
	gl.GetProgramiv(programId, gl.LINK_STATUS, &glErr)
	if glErr == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(programId, gl.INFO_LOG_LENGTH, &logLength)
		log := string(make([]byte, logLength+1))
		gl.GetProgramInfoLog(programId, logLength, nil, gl.Str(log))
		return errors.New("Error while linking program: " + log)
	}
	return nil
}

func LoadShader(shaderPath string, shaderType uint32) (uint32, error) {
	shaderCode, err := os.ReadFile(shaderPath)
	if err != nil {
		return 0, err
	}

	return CreateShader(string(shaderCode), shaderType)
}

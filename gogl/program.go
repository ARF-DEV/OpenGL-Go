package gogl

// import (
// 	"errors"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/go-gl/gl/v3.3-core/gl"
// )

// type Shader struct {
// 	ID           uint32
// 	shaderType   uint32
// 	path         string
// 	lastModified time.Time
// }

// func CreateShader(path string, shaderType uint32) (*Shader, error) {
// 	shaderCode, err := os.ReadFile(path)

// 	if err != nil {
// 		return nil, err
// 	}

// 	shaderId := gl.CreateShader(shaderType)
// 	cSource, free := gl.Strs(string(shaderCode) + "\x00")
// 	gl.ShaderSource(shaderId, 1, cSource, nil)
// 	defer free()
// 	gl.CompileShader(shaderId)
// 	var glErr int32
// 	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &glErr)
// 	if glErr == gl.FALSE {
// 		var logLength int32
// 		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLength)
// 		log := string(make([]byte, logLength+1))
// 		gl.GetShaderInfoLog(shaderId, logLength, nil, gl.Str(log))
// 		return nil, errors.New("Error Failed to compile shader:\n" + log)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	fileInfo, _ := os.Stat(path)

// 	return &Shader{
// 		ID:           shaderId,
// 		shaderType:   shaderType,
// 		path:         path,
// 		lastModified: fileInfo.ModTime(),
// 	}, nil
// }

// func (s *Shader) IsUpdated() bool {
// 	fInfo, _ := os.Stat(s.path)
// 	return fInfo.ModTime().After(s.lastModified)
// }

// type Program struct {
// 	ID      uint32
// 	shaders []*Shader
// }

// func CreateProgram(shaders ...*Shader) (Program, error) {
// 	result := Program{}
// 	pId := gl.CreateProgram()
// 	for _, shader := range shaders {
// 		gl.AttachShader(pId, shader.ID)
// 		result.shaders = append(result.shaders, shader)
// 	}

// 	gl.LinkProgram(pId)
// 	var status int32
// 	gl.GetProgramiv(pId, gl.LINK_STATUS, &status)
// 	if status == gl.FALSE {
// 		var logLength int32
// 		gl.GetProgramiv(pId, gl.INFO_LOG_LENGTH, &logLength)
// 		log := strings.Repeat("\x00", int(logLength)+1)
// 		gl.GetProgramInfoLog(pId, logLength, nil, gl.Str(log))
// 		return Program{}, errors.New(log)
// 	}

// 	for _, shader := range shaders {
// 		gl.DeleteShader(shader.ID)
// 	}

// 	result.ID = pId
// 	return result, nil
// }

// func (p *Program) ReloadProgram() {
// 	shaders := []*Shader{}
// 	pId := gl.CreateProgram()
// 	for _, v := range p.shaders {
// 		v, err := CreateShader(v.path, v.shaderType)
// 		if err != nil {
// 			panic(err)
// 		}
// 		shaders = append(shaders, v)
// 		gl.AttachShader(pId, v.ID)
// 	}
// 	gl.LinkProgram(pId)
// 	var status int32
// 	gl.GetProgramiv(pId, gl.LINK_STATUS, &status)
// 	if status == gl.FALSE {
// 		var logLength int32
// 		gl.GetProgramiv(pId, gl.INFO_LOG_LENGTH, &logLength)
// 		log := strings.Repeat("\x00", int(logLength)+1)
// 		gl.GetProgramInfoLog(pId, logLength, nil, gl.Str(log))
// 		panic(log)
// 		// return Program{}, errors.New(log)
// 	}
// 	p.ID = pId
// 	p.shaders = shaders
// }

// func (p *Program) IsUpdated() bool {
// 	for _, shader := range p.shaders {
// 		if shader.IsUpdated() {
// 			return true
// 		}
// 	}
// 	return false
// }

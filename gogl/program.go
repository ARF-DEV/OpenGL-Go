package gogl

import (
	"os"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Shader struct {
	ID           uint32
	shaderType   uint32
	path         string
	lastCompiled time.Time
}

func CreateShaderStruct(path string, shaderType uint32) (Shader, error) {
	sId, err := LoadShader(path, shaderType)
	if err != nil {
		return Shader{}, err
	}

	return Shader{
		ID:           sId,
		shaderType:   shaderType,
		path:         path,
		lastCompiled: time.Time{},
	}, nil
}

func (s *Shader) IsUpdated() bool {
	fInfo, _ := os.Stat(s.path)
	return fInfo.ModTime().After(s.lastCompiled)
}

const (
	VERTEX_SHADER = iota
	FRAGMENT_SHADER
)

type Program struct {
	ID      uint32
	shaders map[uint32]*Shader
}

func getShaderIds(program Program) []uint32 {
	ids := []uint32{}
	for _, shader := range program.shaders {
		ids = append(ids, shader.ID)
	}
	return ids
}

// TODO Implement using ... Operator
func CreateProgramStructFromPaths(shaderTypes []uint32, shaderPaths ...string) (Program, error) {
	result := Program{
		shaders: make(map[uint32]*Shader),
	}
	for idx, shaderType := range shaderTypes {
		s, err := CreateShaderStruct(shaderPaths[idx], shaderType)
		if err != nil {
			return Program{}, err
		}
		result.shaders[shaderType] = &s
	}
	ids := getShaderIds(result)

	pId, err := CreateProgram(ids...)
	if err != nil {
		return Program{}, err
	}

	result.ID = pId
	return result, nil
}

func (p *Program) ReloadProgram() {

	gl.DeleteProgram(p.ID)

	*p, _ = CreateProgramStructFromPaths([]uint32{gl.VERTEX_SHADER, gl.FRAGMENT_SHADER},
		p.shaders[gl.VERTEX_SHADER].path, p.shaders[gl.FRAGMENT_SHADER].path)
}

func (p *Program) IsUpdated() bool {
	for _, shader := range p.shaders {
		if shader.IsUpdated() {
			return true
		}
	}
	return false
}

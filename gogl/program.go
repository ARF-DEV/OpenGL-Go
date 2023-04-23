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

type Program struct {
	ID      uint32
	shaders []*Shader
}

func getShaderIds(program Program) []uint32 {
	ids := []uint32{}
	for _, shader := range program.shaders {
		ids = append(ids, shader.ID)
	}
	return ids
}

func CreateProgramStruct(shaders ...*Shader) (Program, error) {
	result := Program{}
	for _, shader := range shaders {
		s, err := CreateShaderStruct(shader.path, shader.shaderType)
		if err != nil {
			return Program{}, err
		}
		result.shaders = append(result.shaders, &s)
	}
	ids := getShaderIds(result)

	pId, err := CreateProgram(ids...)
	if err != nil {
		return Program{}, err
	}

	result.ID = pId
	return result, nil
}

func CreateProgramStructFromPaths(shaderTypes []uint32, shaderPaths ...string) (Program, error) {
	result := Program{}
	for idx, shaderType := range shaderTypes {
		s, err := CreateShaderStruct(shaderPaths[idx], shaderType)
		if err != nil {
			return Program{}, err
		}
		result.shaders = append(result.shaders, &s)
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

	*p, _ = CreateProgramStruct(p.shaders...)
}

func (p *Program) IsUpdated() bool {
	for _, shader := range p.shaders {
		if shader.IsUpdated() {
			return true
		}
	}
	return false
}

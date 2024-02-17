package gogl

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Shader struct {
	ID                   uint32
	vertexShader         string
	fragmentShader       string
	vertexLastModified   time.Time
	fragmentLastModified time.Time
}

func (s *Shader) IsUpdated() bool {
	vertexModTime, _ := getFileModifiedTime(s.vertexShader)
	fragModTime, _ := getFileModifiedTime(s.fragmentShader)

	return fragModTime.After(s.fragmentLastModified) || vertexModTime.After(s.vertexLastModified)
}

func (s *Shader) ReloadOnUpdate() {

	if !s.IsUpdated() {
		return
	}

	log.Println("Updated")
	count := 0
	for {

		vSize, _ := getFileSize(s.vertexShader)
		fSize, _ := getFileSize(s.fragmentShader)

		if vSize != 0 && fSize != 0 {
			break
		}
		count++
	}
	log.Printf("%d Loops\n", count)

	new_shader, err := CreateShader(s.vertexShader, s.fragmentShader)
	if err != nil {
		panic(err)
	}
	gl.DeleteProgram(s.ID)
	s.ID = new_shader.ID
	s.vertexShader = new_shader.vertexShader
	s.fragmentShader = new_shader.fragmentShader
	s.vertexLastModified = new_shader.vertexLastModified
	s.fragmentLastModified = new_shader.fragmentLastModified
}

func ReadSource(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	result := ""
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		text := fileScanner.Text()
		fmt.Println(fileScanner.Bytes())
		fmt.Println(fileScanner.Text())
		result += text + "\n"
	}

	return result, nil

}
func CreateShader(vertexShaderPath string, fragmentShaderPath string) (*Shader, error) {
	vertexId := gl.CreateShader(gl.VERTEX_SHADER)

	vSource, err := os.ReadFile(vertexShaderPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// fmt.Println(string(vSource))
	cStr, free := gl.Strs(string(vSource) + "\x00")
	// fmt.Println(string(**cStr))
	gl.ShaderSource(vertexId, 1, cStr, nil)
	gl.CompileShader(vertexId)
	if getShaderParam(vertexId, gl.COMPILE_STATUS) == gl.FALSE {
		logLength := getShaderParam(vertexId, gl.INFO_LOG_LENGTH)
		log := strings.Repeat("\x00", int(logLength)+1)
		gl.GetShaderInfoLog(vertexId, logLength, nil, gl.Str(log))
		return nil, errors.New("Failed while compiling vertex shader: " + log)
	}
	free()

	fragmentId := gl.CreateShader(gl.FRAGMENT_SHADER)
	fSource, err := os.ReadFile(fragmentShaderPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	cStr, free = gl.Strs(string(fSource) + "\x00")
	gl.ShaderSource(fragmentId, 1, cStr, nil)
	gl.CompileShader(fragmentId)
	if getShaderParam(fragmentId, gl.COMPILE_STATUS) == gl.FALSE {
		logLength := getShaderParam(fragmentId, gl.INFO_LOG_LENGTH)
		log := strings.Repeat("\x00", int(logLength)+1)
		gl.GetShaderInfoLog(fragmentId, logLength, nil, gl.Str(log))
		return nil, errors.New("Failed while compiling fragment shader: " + log)
	}
	free()
	for e := gl.GetError(); e != 0; {
		print("Error: ", e)
	}
	programId := gl.CreateProgram()
	gl.AttachShader(programId, vertexId)
	gl.AttachShader(programId, fragmentId)
	gl.LinkProgram(programId)
	if getProgramParam(programId, gl.LINK_STATUS) == gl.FALSE {
		logLength := getProgramParam(programId, gl.INFO_LOG_LENGTH)
		log := strings.Repeat("\x00", int(logLength)+1)
		gl.GetProgramInfoLog(programId, logLength, nil, gl.Str(log))
		return nil, errors.New("Failed while linking program: " + log)
	}

	vertexModTime, _ := getFileModifiedTime(vertexShaderPath)
	fragmentModTime, _ := getFileModifiedTime(fragmentShaderPath)

	return &Shader{
		ID:                   programId,
		vertexShader:         vertexShaderPath,
		fragmentShader:       fragmentShaderPath,
		vertexLastModified:   vertexModTime,
		fragmentLastModified: fragmentModTime,
	}, nil
}

func getFileSize(file string) (int64, error) {
	fInfo, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return fInfo.Size(), nil
}

// func getFileInfo(file string) (os.FileInfo, error) {
// 	fInfo, err := os.Stat(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return fInfo, nil
// }

func getFileModifiedTime(file string) (time.Time, error) {
	fInfo, err := os.Stat(file)
	if err != nil {
		return time.Time{}, err
	}

	return fInfo.ModTime(), nil
}

func getProgramParam(programId uint32, pName uint32) int32 {
	var res int32
	gl.GetProgramiv(programId, pName, &res)
	return res
}

func getShaderParam(shaderId uint32, pName uint32) int32 {
	var res int32
	gl.GetShaderiv(shaderId, pName, &res)
	return res
}
func (s *Shader) SetUniformFloat(name string, value float32) {
	gl.Uniform1f(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), value)
}

func (s *Shader) SetUniformInt(name string, value int32) {
	gl.Uniform1i(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), value)
}

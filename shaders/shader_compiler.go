package shaders

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"strings"
)

const SHADER_DEFAULT_V = `
    #version 420
	#extension GL_ARB_explicit_uniform_location : enable

    layout(location = 0) in vec3 vp;
    void main() {
        gl_Position = vec4(vp, 1);
    }
` + "\x00"

const SHADER_DEFAULT_F = `
    #version 420
	#extension GL_ARB_explicit_uniform_location : enable

    layout (location = 0) out vec4 frag_colour;
    void main() {
        frag_colour = vec4(0.86, 0.76, 0.81, 1.0);
    }
` + "\x00"

type Shader struct {
	handle uint32
}
type Program struct {
	handle  uint32
	shaders []*Shader
}

func Shader_compiler(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("Failed to compile shader %v: %v", source, log)
	}
	return shader, nil
}

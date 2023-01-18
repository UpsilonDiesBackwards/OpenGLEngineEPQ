package shaders

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"strings"
)

const SHADER_DEFAULT_V = `
    #version 420
	#extension GL_ARB_explicit_uniform_location : enable
	#extension GL_ARB_enhanced_layouts : enable

    layout(location = 0) in vec3 vp;

	layout(binding = 1) uniform PerspectiveBlock {
		mat4 project;
		mat4 camera;
		mat4 world;
	};

    void main() {
        gl_Position = project * camera * world * vec4(vp, 1);
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

func CompilerCompiler(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType) // Create new shader

	csources, free := gl.Strs(source) // Open source file in a way that OpenGL likes.
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	// Handle and output any errors
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

func (prog *Program) GetUniformLocation(name string) int32 {
	return gl.GetUniformLocation(prog.handle, gl.Str(name+"\x00"))
}

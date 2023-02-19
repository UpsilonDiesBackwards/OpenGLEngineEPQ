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

const SHADER_PHONG_V = `
#version 410 core

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;

uniform mat4 model;
uniform mat4 view;
uniform mat4 project;

uniform vec3 lightPos;  // only need one light for a basic example

out vec3 Normal;
out vec3 FragPos;
out vec3 LightPos;

void main()
{
    gl_Position = project * view * model * vec4(position, 1.0);

    // transform positions and vectors to view space before performing lighting
    FragPos = vec3(view * model * vec4(position, 1.0));

    // LightPos = vec3(view * vec4(lightPos, 1.0));
    LightPos = vec3(view * vec4(lightPos, 1.0));

    // transform the normals to the view space
    mat3 normMatrix = mat3(transpose(inverse(view))) * mat3(transpose(inverse(model)));
    Normal = normMatrix * normal;
}
`

const SHADER_PHONG_F = `
#version 410 core

in vec3 Normal;
in vec3 FragPos;
in vec3 LightPos;
out vec4 color;

uniform vec3 objectColor;
uniform vec3 lightColor;

void main()
{
	// affects diffuse and specular lighting
	float lightPower = 2.0f;

	// diffuse and specular intensity are affected by the amount of light they get based on how
	// far they are from a light source (inverse square of distance)
	float distToLight = length(LightPos - FragPos);

	// this is not the correct equation for light decay but it is close
	// see light-casters sample for the proper way
	float distIntensityDecay = 1.0f / pow(distToLight, 2);

	float ambientStrength = 0.05f;
	vec3 ambientLight = ambientStrength * lightColor;

	vec3 norm = normalize(Normal);
	vec3 dirToLight = normalize(LightPos - FragPos);
	float lightNormalDiff = max(dot(norm, dirToLight), 0.0);

	// diffuse light is greatest when surface is perpendicular to light (dot product)
	vec3 diffuse = lightNormalDiff * lightColor;
	vec3 diffuseLight = lightPower * diffuse * distIntensityDecay * lightColor;

	float specularStrength = 1.0f;
	int shininess = 64;
	vec3 viewPos = vec3(0.0f, 0.0f, 0.0f);
	vec3 dirToView = normalize(viewPos - FragPos);
	vec3 reflectDir = reflect(-dirToLight, norm);
	float spec = pow(max(dot(dirToView, reflectDir), 0.0), shininess);
	vec3 specularLight = lightPower * specularStrength * spec * distIntensityDecay * lightColor;

	vec3 result = (diffuseLight + specularLight + ambientLight) * objectColor;
	color = vec4(result, 1.0f);
}
`

type Shader struct {
	handle uint32
}
type Program struct {
	handle  uint32
	shaders []*Shader
}

func ShaderCompiler(source string, shaderType uint32) (uint32, error) {
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

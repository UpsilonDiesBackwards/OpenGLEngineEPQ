#version 420
#extension GL_ARB_explicit_uniform_location : enable
#extension GL_ARB_enhanced_layouts : enable

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;

layout(binding = 1) uniform PerspectiveBlock {
    mat4 project;
    mat4 camera;
    mat4 world;
};

layout(location=0) uniform vec3 lightPos;  // only need one light for a basic example

layout(location = 1) out vec3 Normal;
layout(location = 2) out vec3 FragPos;
layout(location = 3) out vec3 LightPos;

void main()
{
    gl_Position = project * camera * world * vec4(position, 1.0);

    // we transform positions and vectors to view space before performing lighting
    // calculations in the fragment shader so that we know that the viewer position is (0,0,0)
    // FragPos = vec3(world * vec4(position, 1.0));
    FragPos = vec3(camera * world * vec4(position, 1.0));

    // LightPos = vec3(camera * vec4(lightPos, 1.0));
    LightPos = vec3(camera * vec4(lightPos, 1.0));

    // transform the normals to the view space
    mat3 normMatrix = mat3(transpose(inverse(camera))) * mat3(transpose(inverse(world)));
    Normal = normMatrix * normal;
}
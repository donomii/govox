package govox

var (
	vertexShader = `
#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec4 colour;

out vec4 pointcol;

void main() {
	gl_Position = projection * camera * model * vec4(vert, 1);
	pointcol = colour;
}
` + "\x00"

	fragmentShader = `
#version 330

uniform vec4 col;
in vec4 pointcol;
out vec4 outputColor;

void main() {
	outputColor = pointcol;
}
` + "\x00"
)

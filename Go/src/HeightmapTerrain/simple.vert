#version 430

layout (location = 0) in vec3 vertPos;
layout (location = 1) in vec3 vertNormal;
layout (location = 2) in vec2 vertUV;

out vec3 vPos;
out vec3 vNormal;
out vec2 vUV;

void main() {
    vPos = vertPos;
    vNormal = normalize(vertNormal);
    vUV = vertUV;
}

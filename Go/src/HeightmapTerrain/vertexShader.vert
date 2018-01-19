#version 430

layout (location = 0) in vec3 vertPos;
layout (location = 1) in vec3 vertNormal;

// Normal attributes
uniform mat4 viewProjectionMat;
uniform mat4 modelMat;

out vec3 normal;
out vec3 pos;

void main() {

    normal = normalize(vertNormal);
    pos = (modelMat * vec4(vertPos,1)).xyz;

    gl_Position = viewProjectionMat * modelMat * (vec4(vertPos,1));

}


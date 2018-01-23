#version 430

layout (location = 0) in vec3 vertPos;
layout (location = 1) in vec3 vertNormal;
layout (location = 2) in vec2 vertUV;

// Normal attributes
uniform mat4 viewProjectionMat;
uniform mat4 modelMat;

out vec3 normal;
out vec3 pos;
out vec2 uv;
out vec4 screenPos;

void main() {


    // No normal matrix right now, because we don't have any rotation right now anyway.
    normal = normalize(vertNormal);


    vec4 tmpPos = modelMat * vec4(vertPos,1);
    pos = tmpPos.xyz;

    uv = vertUV;

    screenPos = viewProjectionMat * tmpPos;
    gl_Position = screenPos;
    screenPos.xyz /= screenPos.w;

}


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

void main() {



    normal = normalize(vertNormal);


    vec4 tmpPos = modelMat * vec4(vertPos,1);
    pos = tmpPos.xyz;

    uv = vertUV;

    gl_Position = viewProjectionMat * tmpPos;

}


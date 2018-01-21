#version 430

layout (location = 0) in vec3 vertPos;
layout (location = 1) in vec3 vertNormal;
layout (location = 2) in vec2 vertUV;

// Normal attributes
uniform mat4 viewProjectionMat;
uniform mat4 modelMat;

uniform sampler2D heightmapTextureOriginal;
uniform sampler2D heightmapTextureMerged;
uniform sampler2D heightmapTexture900m;
uniform vec2 textureSize;

out vec3 normal;
out vec3 pos;
out vec2 uv;

vec3 calcOffsetPosition(vec3 pos, vec2 offset) {
    float h = texture(heightmapTextureMerged, vertUV + offset * vec2( 1.0/textureSize.x, 1.0/textureSize.y)).r;
    return pos + vec3(offset.x, h * 0.4 - 0.2, offset.y);
}

void main() {



    normal = normalize(vertNormal);

    vec3 horizontal = calcOffsetPosition(vertPos, vec2(1,0)) - calcOffsetPosition(vertPos, vec2(-1,0));
    vec3 vertical   = calcOffsetPosition(vertPos, vec2(0,1)) - calcOffsetPosition(vertPos, vec2(0,-1));

    //normal = normalize(cross(vertical, horizontal));


    float height = texture(heightmapTextureOriginal, vertUV).r;
    vec4 tmpPos = modelMat * vec4(vertPos + vec3(0,height * 0.6 - 0.3,0),1);
    //vec4 tmpPos = modelMat * vec4(vertPos,1);
    pos = tmpPos.xyz;

    uv = vertUV;

    gl_Position = viewProjectionMat * tmpPos;

}


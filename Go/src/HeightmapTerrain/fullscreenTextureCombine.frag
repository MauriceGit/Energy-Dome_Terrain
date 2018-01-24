#version 430

uniform highp sampler2D terrainTexture;
uniform highp sampler2D sphereTexture;

in vec2 fUV;
out vec4 fragColor;

void main()
{
    vec4 sphereColor = texture(sphereTexture, fUV);

    vec2 newUV = fUV;
    newUV.x += sphereColor.r*0.2;

    fragColor = texture(terrainTexture, newUV) + sphereColor;




}


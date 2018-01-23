#version 430

uniform highp sampler2D overwriteTexture;

in vec2 fUV;
out vec4 fragColor;

void main()
{
    fragColor = texture(overwriteTexture, fUV);
}


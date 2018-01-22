#version 430

in vec3 pos;
in vec3 normal;
in vec2 uv;

uniform vec3 color;
uniform float dt;

uniform sampler2D energyTexture;
uniform sampler2D energyAnimationTexture;
uniform int polygonMode;

out vec4 colorOut;

void main() {


    float animation = texture(energyAnimationTexture, uv).r * 2. -1.;

    vec2 newUV   = uv + vec2(animation*0.1,dt*0.13 + animation*0.05);
    float energy = texture(energyTexture, newUV).r;

    energy *= 1.2;
    energy  = pow(energy, 3);
    energy *= 5.0;

    colorOut = energy * vec4(color,1);





}



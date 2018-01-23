#version 430

in vec3 pos;
in vec3 normal;
in vec2 uv;

uniform vec3 color;
uniform vec3 light;
uniform bool isLight;

uniform sampler2D heightmapTextureOriginal;
uniform sampler2D heightmapTextureMerged;
uniform sampler2D heightmapTexture900m;
uniform vec2 textureSize;

uniform int polygonMode;

out vec4 colorOut;

void main() {
    colorOut = vec4(color, 1);




    vec3 normal2 = normal;




    vec3 l = normalize(light - pos);

    vec3 specularColor = color*3;
    float dotProduct = max(dot(normal2,l), 0.0);
    vec3 specular = specularColor * pow(dotProduct, 8.0);
    specular = clamp(specular, 0.0, 1.0);

    vec3 diffuseColor = color*2;
    vec3 diffuse  = diffuseColor * max(dot(normal2, l), 0);
    diffuse = clamp(diffuse, 0.0, 1.0);

    vec3 diffuseColorNeg = color*3;
    vec3 diffuseNeg  = diffuseColorNeg * max(dot(-normal2, l), 0);
    diffuseNeg = clamp(diffuseNeg, 0.0, 1.0);
    diffuseNeg = vec3(1)-diffuseNeg;

    vec3 ambient = color / 1.5;

    if (isLight) {
        colorOut = vec4(color, 1.0);
    } else {
        colorOut = vec4(diffuseNeg/4 + diffuse/4 + ambient/4 + specular/4 + color/3, 1.0);
    }

    //colorOut = vec4(0.3,0.08,0.08,0.2);
    vec4 texColor = texture(heightmapTextureOriginal, uv);
    colorOut = mix(vec4(1.,0.05,0.05,1), vec4(0.15,0.8,0.15,1), texColor.r*3.);

    //if (polygonMode == 0) {
    //    colorOut *= 0.6;
    //} else if (polygonMode == 1) {
    //    colorOut *= 0.1;
    //} else if (polygonMode == 2) {
    //    colorOut *= 0.15;
    //}

    colorOut.a = 1.0;

}



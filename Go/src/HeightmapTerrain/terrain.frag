#version 430

in vec3 tePos;
in vec3 teNormal;
in vec2 teUV;

uniform vec3 color;
uniform vec3 light;
uniform bool isLight;

uniform sampler2D heightmapTextureOriginal;
uniform sampler2D heightmapTextureMerged;
uniform sampler2D heightmapTexture900m;
uniform sampler2D terrainNormalHeightTexture;
uniform vec2 textureSize;

uniform int polygonMode;

out vec4 colorOut;

vec4 calcHeightColor() {
    vec4 cOut;
    float texColor = texture(heightmapTextureOriginal, teUV).r * 4.;
    texColor = pow(texColor,2);

    vec4 c0 = vec4(1,0,0,1);
    vec4 c1 = vec4(0.15,0.6,0.12,1);
    vec4 c2 = vec4(0.8,0.8,0.2,1);

    if (texColor <= 0.5) {
        cOut = mix(c0, c1, texColor*2.);
    } else {
        cOut = mix(c1, c2, (texColor-0.5)*2.);
    }
    return cOut;
}

void main() {


    colorOut = calcHeightColor();
    //colorOut = vec4(color, 1);
    //vec3 heightColor = colorOut.rgb;
    //
    //
    //
    //
    //
    //vec3 normal2 = texture(terrainNormalHeightTexture, teUV).rgb;
    //
    //vec3 l = normalize(light - tePos);
    //
    //vec3 specularColor = heightColor*3;
    //float dotProduct = max(dot(normal2,l), 0.0);
    //vec3 specular = specularColor * pow(dotProduct, 8.0);
    //specular = clamp(specular, 0.0, 1.0);
    //
    //vec3 diffuseColor = heightColor*2;
    //vec3 diffuse  = diffuseColor * max(dot(normal2, l), 0);
    //diffuse = clamp(diffuse, 0.0, 1.0);
    //
    //vec3 diffuseColorNeg = heightColor*3;
    //vec3 diffuseNeg  = diffuseColorNeg * max(dot(-normal2, l), 0);
    //diffuseNeg = clamp(diffuseNeg, 0.0, 1.0);
    //diffuseNeg = vec3(1)-diffuseNeg;
    //
    //vec3 ambient = heightColor / 1.5;
    //
    //if (isLight) {
    //    colorOut = vec4(heightColor, 1.0);
    //} else {
    //    colorOut = vec4(diffuseNeg/4 + diffuse/4 + ambient/4 + specular/4 + heightColor/3, 1.0);
    //}



}

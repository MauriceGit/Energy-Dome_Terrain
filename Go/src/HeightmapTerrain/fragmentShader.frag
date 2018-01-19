#version 430

in vec3 normal;
in vec3 pos;

uniform vec3 color;
uniform vec3 light;
uniform bool isLight;

out vec4 colorOut;

void main() {
    colorOut = vec4(color, 1);

    vec3 l = normalize(light - pos);

    vec3 specularColor = color*3;
    float dotProduct = max(dot(normal,l), 0.0);
    vec3 specular = specularColor * pow(dotProduct, 8.0);
    specular = clamp(specular, 0.0, 1.0);

    vec3 diffuseColor = color*2;
    vec3 diffuse  = diffuseColor * max(dot(normal, l), 0);
    diffuse = clamp(diffuse, 0.0, 1.0);

    vec3 diffuseColorNeg = color*3;
    vec3 diffuseNeg  = diffuseColorNeg * max(dot(-normal, l), 0);
    diffuseNeg = clamp(diffuseNeg, 0.0, 1.0);
    diffuseNeg = vec3(1)-diffuseNeg;

    vec3 ambient = color / 1.5;

    if (isLight) {
        colorOut = vec4(color, 1);
    } else {
        colorOut = vec4(diffuseNeg/4 + diffuse/4 + ambient/4 + specular/4 + color/3, 1.0);
    }

}



#version 430

in vec3 pos;
in vec3 normal;
in vec2 uv;
in vec4 screenPos;

uniform vec3 color;
uniform float dt;

uniform vec3 camPos;

uniform sampler2D energyTexture;
uniform sampler2D energyAnimationTexture;
uniform int polygonMode;

uniform sampler2D sceneColorTex;
uniform sampler2D sceneDepthTex;

out vec4 colorOut;

float linearizeDepth (float depth) {
    float nearPlane = 0.1, farPlane = 2000.0;
    return (2.0*nearPlane) / (farPlane + nearPlane - depth * (farPlane - nearPlane));
}

void main() {

    vec3 sPos = screenPos.xyz;
    sPos = sPos *0.5 +0.5;

    float sphereDepth = screenPos.z*0.5+0.5;
    float sceneDepth  = texture(sceneDepthTex, gl_FragCoord.xy/vec2(1000)).r;

    // Do the depth testing manually because we render this in post-processing.
    if (sphereDepth > sceneDepth) {
        discard;
    }

    float animation = texture(energyAnimationTexture, uv + vec2(dt*0.11, -dt*0.015)).r * 2. -1.;

    vec2 newUV   = uv + vec2(animation*0.1,dt*0.13 + animation*0.05);
    float energy = texture(energyTexture, newUV).r;

    energy *= 1.2;
    energy  = pow(energy, 3);
    energy *= 5.0;

    colorOut = energy * vec4(color,1);

    vec3 viewVec = normalize(camPos - pos);
    float d = dot(viewVec, normal);

    // Have a nice gradient to the defined color at all edges of the sphere.
    if (d > 0 && d < 0.5) {
        colorOut.rgb = mix(colorOut.rgb, color, pow(1.0-d*2.,3));
    }

    float linearSphereDepth = linearizeDepth(sphereDepth);
    float linearSceneDepth  = linearizeDepth(sceneDepth);
    float linearDepthDiff   = abs(linearSceneDepth-linearSphereDepth);

    if (linearDepthDiff <= 0.005) {
        //colorOut.rgb = vec3(0.3,0,0);

        colorOut.rgb = mix(colorOut.rgb, color, pow(1.0-linearDepthDiff*200.,3));

    }



}



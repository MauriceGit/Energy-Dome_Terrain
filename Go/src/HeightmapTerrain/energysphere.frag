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

uniform vec2 nearFarPlane;
uniform vec2 windowSize;

out vec4 colorOut;

float linearizeDepth (float depth) {
    return (2.0*nearFarPlane.x) / (nearFarPlane.y + nearFarPlane.x - depth * (nearFarPlane.y - nearFarPlane.x));
}

void main() {

    vec3 sPos = screenPos.xyz;
    sPos = sPos *0.5 +0.5;

    float sphereDepth = screenPos.z*0.5+0.5;
    float sceneDepth  = texture(sceneDepthTex, gl_FragCoord.xy/windowSize).r;

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
    float sphereEdgeHighlightDistance = 0.5;

    // Have a nice gradient to the defined color at all edges of the sphere.
    if (d > 0 && d < sphereEdgeHighlightDistance) {
        colorOut.rgb = mix(colorOut.rgb, color, pow(1.0-d*(1.0/sphereEdgeHighlightDistance),3));
    }

    float linearSphereDepth = linearizeDepth(sphereDepth);
    float linearSceneDepth  = linearizeDepth(sceneDepth);
    float linearDepthDiff   = abs(linearSceneDepth-linearSphereDepth);
    float geometryHighlightDistance = 0.005;

    // Have a really nice gradient on all geometry defined edges! This is looking much better as expected :)
    if (linearDepthDiff <= geometryHighlightDistance) {
        colorOut.rgb = mix(colorOut.rgb, color, pow(1.0-linearDepthDiff*(1.0/geometryHighlightDistance),3));
    }

}



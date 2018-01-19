package camera

import (
    "github.com/go-gl/mathgl/mgl32"
    "github.com/go-gl/glfw/v3.2/glfw"
)

const (
    PI = 3.14159265359
    ROTATE_SCALE   = 0.5
    UPDOWN_SCALE   = 0.15
    DISTANCE_SCALE = 0.1
)

var g_cameraPos mgl32.Vec3 = mgl32.Vec3{0,8,15}
var g_center    mgl32.Vec3 = mgl32.Vec3{0,0,0}
var g_up        mgl32.Vec3 = mgl32.Vec3{0,1,0}

var g_angle float32 = 0.0
var g_cursorPos mgl32.Vec2
var g_leftButtonDown bool  = false



func toRad(a float32) float32 {
    return a*PI/180.
}

// We ignore horizontal scrolling.
func UpdateMouseScroll(xpos, ypos float64) {
    var scale float32 = 1.0
    switch {
        case ypos > 0:
            scale -= DISTANCE_SCALE
        case ypos < 0:
            scale += DISTANCE_SCALE
    }

    g_cameraPos = g_cameraPos.Mul(scale)
}

// We only work with the left mouse button and ignore modifiers for now!
func UpdateMouseButton(button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
    g_leftButtonDown = button == glfw.MouseButtonLeft && action == glfw.Press
}

func UpdateCursorPos(xpos, ypos float64) {
    new  := mgl32.Vec2{float32(xpos), float32(ypos)}
    diff := mgl32.Vec2{float32(xpos), float32(ypos)}.Sub(g_cursorPos)
    g_cursorPos = new

    if g_leftButtonDown {

        r := g_cameraPos.Len()
        rotateY := mgl32.Rotate3DY(toRad(-diff.X())*ROTATE_SCALE)
        g_cameraPos = rotateY.Mul3x1(g_cameraPos)
        g_cameraPos = g_cameraPos.Add(mgl32.Vec3{0,diff.Y()*UPDOWN_SCALE,0}).Normalize().Mul(r)

    }
}

func GetCameraLookAt() (mgl32.Vec3, mgl32.Vec3, mgl32.Vec3) {
    return g_cameraPos, g_center, g_up
}

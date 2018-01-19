package main

import (
    . "HeightmapTerrain/Geometry"
    . "HeightmapTerrain/Camera"
    . "HeightmapTerrain/OpenGL"
    "runtime"
    "github.com/go-gl/mathgl/mgl32"
    "fmt"
    "github.com/go-gl/gl/v4.5-core/gl"
    "github.com/go-gl/glfw/v3.2/glfw"
)

// Constants and global variables

const (
    g_WindowWidth  = 1000
    g_WindowHeight = 1000

    g_cubeWidth    = 10
    g_cubeHeight   = 10
    g_cubeDepth    = 10

)

const g_WindowTitle  = "Heightmap Terrain"
var g_ShaderID uint32


// Normal Camera
var g_fovy      = mgl32.DegToRad(90.0)
var g_aspect    = float32(g_WindowWidth)/g_WindowHeight
var g_nearPlane = float32(0.1)
var g_farPlane  = float32(2000.0)

var g_viewMatrix          mgl32.Mat4

var g_light Object


var g_timeSum float32 = 0.0
var g_lastCallTime float64 = 0.0
var g_frameCount int = 0
var g_fps float32 = 60.0

var g_fillMode = 0

func init() {
    // GLFW event handling must run on the main OS thread
    runtime.LockOSThread()
}


func printHelp() {
    fmt.Println(
        `Help yourself.`,
    )
}

// Set OpenGL version, profile and compatibility
func initGraphicContext() (*glfw.Window, error) {
    glfw.WindowHint(glfw.Resizable, glfw.True)
    glfw.WindowHint(glfw.ContextVersionMajor, 4)
    glfw.WindowHint(glfw.ContextVersionMinor, 3)
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
    glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

    window, err := glfw.CreateWindow(g_WindowWidth, g_WindowHeight, g_WindowTitle, nil, nil)
    if err != nil {
        return nil, err
    }
    window.MakeContextCurrent()

    // Initialize Glow
    if err := gl.Init(); err != nil {
        return nil, err
    }

    return window, nil
}

func defineModelMatrix(shader uint32, pos, scale mgl32.Vec3) {
    matScale := mgl32.Scale3D(scale.X(), scale.Y(), scale.Z())
    matTrans := mgl32.Translate3D(pos.X(), pos.Y(), pos.Z())
    model := matTrans.Mul4(matScale)
    modelUniform := gl.GetUniformLocation(shader, gl.Str("modelMat\x00"))
    gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
}

// Defines the Model-View-Projection matrices for the shader.
func defineMatrices(shader uint32) {
    projection := mgl32.Perspective(g_fovy, g_aspect, g_nearPlane, g_farPlane)
    camera := mgl32.LookAtV(GetCameraLookAt())

    viewProjection := projection.Mul4(camera);
    cameraUniform := gl.GetUniformLocation(shader, gl.Str("viewProjectionMat\x00"))
    gl.UniformMatrix4fv(cameraUniform, 1, false, &viewProjection[0])
}

func renderObject(shader uint32, obj Object) {

    // Model transformations are now encoded per object directly before rendering it!
    defineModelMatrix(shader, obj.Pos, obj.Scale)

    gl.BindVertexArray(obj.Geo.VertexObject)

    gl.Uniform3fv(gl.GetUniformLocation(shader, gl.Str("color\x00")), 1, &obj.Color[0])
    gl.Uniform3fv(gl.GetUniformLocation(shader, gl.Str("light\x00")), 1, &g_light.Pos[0])
    var isLighti int32 = 0
    if obj.IsLight {
        isLighti = 1
    }
    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("isLight\x00")), isLighti)

    gl.DrawArrays(gl.TRIANGLES, 0, obj.Geo.VertexCount)

    gl.BindVertexArray(0)

}

func renderEverything(shader uint32) {

    gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
    gl.Enable(gl.DEPTH_TEST)
    // Nice blueish background
    gl.ClearColor(135.0/255.,206.0/255.,235.0/255., 1.0)

    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.Viewport(0, 0, g_WindowWidth, g_WindowHeight)

    gl.UseProgram(shader)

    defineMatrices(shader)
    renderObject(shader, g_light)

    gl.UseProgram(0)

}

// Callback method for a keyboard press
func cbKeyboard(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

    // All changes come VERY easy now.
    if action == glfw.Press {
        switch key {
            // Close the Simulation.
            case glfw.KeyEscape, glfw.KeyQ:
                window.SetShouldClose(true)
            case glfw.KeyH:
                printHelp()
            case glfw.KeySpace:
            case glfw.KeyF1:
                g_fillMode += 1
                switch (g_fillMode%3) {
                    case 0:
                        gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
                    case 1:
                        gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
                    case 2:
                        gl.PolygonMode(gl.FRONT_AND_BACK, gl.POINT)
                }
            case glfw.KeyF2:

            case glfw.KeyF3:
            case glfw.KeyUp:
                g_light.Pos = g_light.Pos.Add(mgl32.Vec3{0,1.0,0})
            case glfw.KeyDown:
                g_light.Pos = g_light.Pos.Add(mgl32.Vec3{0,-1.0,0})
            case glfw.KeyLeft:
            case glfw.KeyRight:
        }
    }

}

// see: https://github.com/go-gl/glfw/blob/master/v3.2/glfw/input.go
func cbMouseScroll(window *glfw.Window, xpos, ypos float64) {
    UpdateMouseScroll(xpos, ypos)
}

func cbMouseButton(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
    UpdateMouseButton(button, action, mods)
}

func cbCursorPos(window *glfw.Window, xpos, ypos float64) {
    UpdateCursorPos(xpos, ypos)
}


// Register all needed callbacks
func registerCallBacks (window *glfw.Window) {
    window.SetKeyCallback(cbKeyboard)
    window.SetScrollCallback(cbMouseScroll)
    window.SetMouseButtonCallback(cbMouseButton)
    window.SetCursorPosCallback(cbCursorPos)
}


func displayFPS(window *glfw.Window) {
    currentTime := glfw.GetTime()
    g_timeSum += float32(currentTime - g_lastCallTime)


    if g_frameCount%60 == 0 {
        g_fps = float32(1.0) / (g_timeSum/60.0)
        g_timeSum = 0.0

        s := fmt.Sprintf("FPS: %.1f", g_fps)
        window.SetTitle(s)
    }

    g_lastCallTime = currentTime
    g_frameCount += 1

}

// Mainloop for graphics updates and object animation
func mainLoop (window *glfw.Window) {

    registerCallBacks(window)
    glfw.SwapInterval(0)

    for !window.ShouldClose() {

        displayFPS(window)

        // This actually renders everything.
        renderEverything(g_ShaderID)

        window.SwapBuffers()
        glfw.PollEvents()
    }



}

func main() {
    var err error = nil
    if err = glfw.Init(); err != nil {
        panic(err)
    }
    // Terminate as soon, as this the function is finished.
    defer glfw.Terminate()

    window, err := initGraphicContext()
    if err != nil {
        // Decision to panic or do something different is taken in the main
        // method and not in sub-functions
        panic(err)
    }

    path := "../Go/src/HeightmapTerrain/"
    g_ShaderID, err = NewProgram(path+"vertexShader.vert", path+"fragmentShader.frag")
    if err != nil {
        panic(err)
    }

    g_light = CreateObject(CreateUnitSphere(10), mgl32.Vec3{3,15,0}, mgl32.Vec3{0.2,0.2,0.2}, mgl32.Vec3{1,1,0}, true)

    mainLoop(window)

}





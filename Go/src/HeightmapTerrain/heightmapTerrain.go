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
    g_windowWidth  = 1000
    g_windowHeight = 1000

    g_cubeWidth    = 10
    g_cubeHeight   = 10
    g_cubeDepth    = 10

)

const g_windowTitle  = "Heightmap Terrain"
var g_terrainShaderID uint32
var g_energySphereShaderID uint32
var g_fullscreenTexturedShaderID uint32

// Framebuffer object with color and depth attachments
var g_terrainFbo uint32
var g_terrainColorTex uint32
var g_terrainDepthTex uint32
var g_sphereFbo uint32
var g_sphereColorTex uint32
var g_sphereDepthTex uint32
// Multisampled version
var g_sceneFboMS uint32
var g_sceneColorTexMS uint32
var g_sceneDepthTexMS uint32

// Multisampling
var g_multisamplingEnabled bool = true

// Tessellation factor
var g_tessellationSubdivision int32 = 7

// Normal Camera
var g_fovy      = mgl32.DegToRad(90.0)
var g_aspect    = float32(g_windowWidth)/g_windowHeight
var g_nearPlane = float32(0.1)
var g_farPlane  = float32(2000.0)

var g_viewMatrix          mgl32.Mat4

//var g_light   Object
var g_terrain Object
var g_sphere  Object
var g_fullscreenQuad Object
var g_heightmapTextureOriginal ImageTexture
var g_heightmapTexture900m     ImageTexture
var g_heightmapTextureMerged   ImageTexture
var g_energyTexture            ImageTexture
var g_energyAnimationTexture   ImageTexture


var g_timeSum float32 = 0.0
var g_currentTime float64 = 0.0
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

    window, err := glfw.CreateWindow(g_windowWidth, g_windowHeight, g_windowTitle, nil, nil)
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

func renderTerrain(shader uint32, obj Object) {

    // Model transformations are now encoded per object directly before rendering it!
    defineModelMatrix(shader, obj.Pos, obj.Scale)

    gl.BindVertexArray(obj.Geo.VertexBuffer)
    gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, obj.Geo.IndexBuffer)

    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("polygonMode\x00")), int32(g_fillMode))

    gl.Uniform3fv(gl.GetUniformLocation(shader, gl.Str("color\x00")), 1, &obj.Color[0])

    lightPos := mgl32.Vec3{60,80,0}

    gl.Uniform3fv(gl.GetUniformLocation(shader, gl.Str("light\x00")), 1, &lightPos[0])
    var isLighti int32 = 0
    if obj.IsLight {
        isLighti = 1
    }
    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("isLight\x00")), isLighti)

    gl.ActiveTexture(gl.TEXTURE0)
    gl.BindTexture(gl.TEXTURE_2D, g_heightmapTextureOriginal.TextureHandle)
    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("heightmapTextureOriginal\x00")), 0)

    gl.ActiveTexture(gl.TEXTURE0+1)
    gl.BindTexture(gl.TEXTURE_2D, g_heightmapTextureMerged.TextureHandle)
    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("heightmapTextureMerged\x00")), 1)

    gl.ActiveTexture(gl.TEXTURE0+2)
    gl.BindTexture(gl.TEXTURE_2D, g_heightmapTexture900m.TextureHandle)
    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("heightmapTexture900m\x00")), 2)

    textureSize := g_heightmapTextureMerged.TextureSize
    gl.Uniform2fv(gl.GetUniformLocation(shader, gl.Str("textureSize\x00")), 1, &textureSize[0])

    camPos, _, _ := GetCameraLookAt()
    gl.Uniform3fv(gl.GetUniformLocation(shader, gl.Str("camPos\x00")), 1, &camPos[0])
    nearFarPlane := mgl32.Vec2{g_nearPlane, g_farPlane}
    gl.Uniform2fv(gl.GetUniformLocation(shader, gl.Str("nearFarPlane\x00")), 1, &nearFarPlane[0])

    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("tessSubdivInner\x00")),  g_tessellationSubdivision)
    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("tessSubdivOuterU\x00")), g_tessellationSubdivision)
    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("tessSubdivOuterV\x00")), g_tessellationSubdivision)
    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("tessSubdivOuterW\x00")), g_tessellationSubdivision)
    gl.PatchParameteri(gl.PATCH_VERTICES, 3)

    gl.DrawElements(gl.PATCHES, obj.Geo.IndexCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
}

func renderEnergySphere(shader uint32, obj Object) {

    // Model transformations are now encoded per object directly before rendering it!
    defineModelMatrix(shader, obj.Pos, obj.Scale)

    gl.BindVertexArray(obj.Geo.VertexBuffer)
    gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, obj.Geo.IndexBuffer)

    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("polygonMode\x00")), int32(g_fillMode))

    gl.Uniform1f(gl.GetUniformLocation(shader, gl.Str("dt\x00")), float32(g_currentTime))

    gl.Uniform3fv(gl.GetUniformLocation(shader, gl.Str("color\x00")), 1, &obj.Color[0])

    camPos, _, _ := GetCameraLookAt()
    gl.Uniform3fv(gl.GetUniformLocation(shader, gl.Str("camPos\x00")), 1, &camPos[0])


    nearFarPlane := mgl32.Vec2{g_nearPlane, g_farPlane}
    gl.Uniform2fv(gl.GetUniformLocation(shader, gl.Str("nearFarPlane\x00")), 1, &nearFarPlane[0])
    windowSize := mgl32.Vec2{g_windowWidth, g_windowHeight}
    gl.Uniform2fv(gl.GetUniformLocation(shader, gl.Str("windowSize\x00")), 1, &windowSize[0])



    gl.ActiveTexture(gl.TEXTURE0)
    gl.BindTexture(gl.TEXTURE_2D, g_energyTexture.TextureHandle)
    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("energyTexture\x00")), 0)
    gl.ActiveTexture(gl.TEXTURE0+1)
    gl.BindTexture(gl.TEXTURE_2D, g_energyAnimationTexture.TextureHandle)
    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("energyAnimationTexture\x00")), 1)
    gl.ActiveTexture(gl.TEXTURE0+2)
    gl.BindTexture(gl.TEXTURE_2D, g_terrainColorTex)
    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("sceneColorTex\x00")), 2)
    gl.ActiveTexture(gl.TEXTURE0+3)
    gl.BindTexture(gl.TEXTURE_2D, g_terrainDepthTex)
    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("sceneDepthTex\x00")), 3)

    gl.DrawElements(gl.TRIANGLES, obj.Geo.IndexCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
}

func renderFullscreenQuad(shader uint32, obj Object) {

    gl.BindVertexArray(obj.Geo.VertexBuffer)
    gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, obj.Geo.IndexBuffer)

    gl.ActiveTexture(gl.TEXTURE0)
    gl.BindTexture(gl.TEXTURE_2D, g_terrainColorTex)
    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("terrainTexture\x00")), 0)
    gl.ActiveTexture(gl.TEXTURE0+1)
    gl.BindTexture(gl.TEXTURE_2D, g_sphereColorTex)
    gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("sphereTexture\x00")), 1)

    gl.DrawElements(gl.TRIANGLES, obj.Geo.IndexCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
}

func renderTextureCombine() {
    gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
    gl.ClearColor(0,0,0,0)
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.Viewport(0, 0, g_windowWidth, g_windowHeight)

    gl.UseProgram(g_fullscreenTexturedShaderID)
    renderFullscreenQuad(g_fullscreenTexturedShaderID, g_fullscreenQuad)
}

func renderEnergySphereFbo() {

    var fbo uint32 = g_sphereFbo
    if g_multisamplingEnabled {
        fbo = g_sceneFboMS
    }

    gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)
    gl.ClearColor(0,0,0,0)
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.Viewport(0, 0, g_windowWidth, g_windowHeight)

    //gl.UseProgram(g_fullscreenTexturedShaderID)
    //renderFullscreenQuad(g_fullscreenTexturedShaderID, g_fullscreenQuad)

    gl.Enable(gl.BLEND)
    gl.Disable(gl.DEPTH_TEST)

    switch (g_fillMode) {
        case 0:
            gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
        case 1:
            gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
        case 2:
            gl.PolygonMode(gl.FRONT_AND_BACK, gl.POINT)
    }

    gl.UseProgram(g_energySphereShaderID)
    defineMatrices(g_energySphereShaderID)
    renderEnergySphere(g_energySphereShaderID, g_sphere)

    gl.Disable(gl.BLEND)
    gl.Enable(gl.DEPTH_TEST)

    gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

    if g_multisamplingEnabled {
        gl.BindFramebuffer(gl.READ_FRAMEBUFFER, g_sceneFboMS)
        gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, g_sphereFbo)
        gl.DrawBuffer(gl.BACK)
        gl.BlitFramebuffer(0, 0, g_windowWidth, g_windowHeight, 0, 0, g_windowWidth, g_windowHeight, gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT, gl.NEAREST)
    }
}

func renderTerrainFbo() {

    var fbo uint32 = g_terrainFbo
    if g_multisamplingEnabled {
        fbo = g_sceneFboMS
    }

    gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)
    gl.ClearColor(0,0,0,0)
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.Viewport(0, 0, g_windowWidth, g_windowHeight)

    switch (g_fillMode) {
        case 0:
            gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
        case 1:
            gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
        case 2:
            gl.PolygonMode(gl.FRONT_AND_BACK, gl.POINT)
    }

    gl.UseProgram(g_terrainShaderID)
    defineMatrices(g_terrainShaderID)
    renderTerrain(g_terrainShaderID, g_terrain)

    gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

    if g_multisamplingEnabled {
        gl.BindFramebuffer(gl.READ_FRAMEBUFFER, g_sceneFboMS)
        gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, g_terrainFbo)
        gl.DrawBuffer(gl.BACK)
        gl.BlitFramebuffer(0, 0, g_windowWidth, g_windowHeight, 0, 0, g_windowWidth, g_windowHeight, gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT, gl.NEAREST)
    }

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
                g_fillMode = (g_fillMode+1) % 3
                //switch (g_fillMode) {
                //    case 0:
                //        gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
                //    case 1:
                //        gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
                //    case 2:
                //        gl.PolygonMode(gl.FRONT_AND_BACK, gl.POINT)
                //}
            case glfw.KeyF2:

            case glfw.KeyF3:
            case glfw.KeyUp:
                g_tessellationSubdivision++
            case glfw.KeyDown:
                if g_tessellationSubdivision > 1 {
                    g_tessellationSubdivision--
                }
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

    g_timeSum += float32(g_currentTime - g_lastCallTime)


    if g_frameCount%60 == 0 {
        g_fps = float32(1.0) / (g_timeSum/60.0)
        g_timeSum = 0.0

        s := fmt.Sprintf("FPS: %.1f", g_fps)
        window.SetTitle(s)
    }

    g_lastCallTime = g_currentTime
    g_frameCount += 1

}

// Mainloop for graphics updates and object animation
func mainLoop (window *glfw.Window) {

    registerCallBacks(window)
    glfw.SwapInterval(0)

    gl.BlendEquation(gl.FUNC_ADD)
    gl.BlendFunc(gl.ONE, gl.ONE)

    for !window.ShouldClose() {

        g_currentTime = glfw.GetTime()
        displayFPS(window)

        renderTerrainFbo()
        renderEnergySphereFbo()
        renderTextureCombine()

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
    g_terrainShaderID, err = NewProgram(path+"simple.vert", path+"simple.tcs", path+"terrain.tes", path+"terrain.frag")
    if err != nil {
        panic(err)
    }
    g_energySphereShaderID, err = NewProgram(path+"energysphere.vert", "", "", path+"energysphere.frag")
    if err != nil {
        panic(err)
    }
    g_fullscreenTexturedShaderID, err = NewProgram(path+"fullscreenTextureCombine.vert", "", "", path+"fullscreenTextureCombine.frag")
    if err != nil {
        panic(err)
    }

    g_terrainFbo   = CreateFbo(&g_terrainColorTex, &g_terrainDepthTex, g_windowWidth, g_windowHeight, false, 1, false, 1)
    g_sphereFbo    = CreateFbo(&g_sphereColorTex, &g_sphereDepthTex, g_windowWidth, g_windowHeight, false, 1, false, 1)
    g_sceneFboMS = CreateFbo(&g_sceneColorTexMS, &g_sceneDepthTexMS, g_windowWidth, g_windowHeight, true, 4, false, 1)

    g_heightmapTextureMerged   = CreateImageTexture(path+"Textures/boeblingen_Height_Map_Merged.png", false)
    g_heightmapTextureOriginal = CreateImageTexture(path+"Textures/boeblingen_Height_Map_Original.png", false)
    g_heightmapTexture900m     = CreateImageTexture(path+"Textures/boeblingen_Height_Map_900m.png", false)

    g_energyTexture            = CreateImageTexture(path+"Textures/tyllo-caustic1_bw_bigger.png", true)
    g_energyAnimationTexture   = CreateImageTexture(path+"Textures/tyllo-caustics02_big.png", true)

    //g_light   = CreateObject(CreateUnitSphere(10), mgl32.Vec3{60,80,0}, mgl32.Vec3{10.2,10.2,10.2}, mgl32.Vec3{0,0,0}, true)

    //energyColor := mgl32.Vec3{0.2,0.55,0.3}
    energyColor := mgl32.Vec3{0.1,0.15,0.55}

    g_terrain = CreateObject(CreateUnitSquareGeometry(10, mgl32.Vec3{0,0,0}), mgl32.Vec3{0,0,0}, mgl32.Vec3{1000.,1000.,1000.}, mgl32.Vec3{139./255.,0,0}, false)
    g_sphere  = CreateObject(CreateUnitSphereGeometry(50, 50), mgl32.Vec3{0,-100,0}, mgl32.Vec3{400.,400.,400.}, energyColor, false)
    g_fullscreenQuad = CreateObject(CreateFullscreenQuadGeometry(), mgl32.Vec3{0,0,0}, mgl32.Vec3{1,1,1}, mgl32.Vec3{0,0,0}, false)

    mainLoop(window)

}





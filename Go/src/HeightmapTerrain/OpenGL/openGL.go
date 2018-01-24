package opengl

import (
    . "HeightmapTerrain/Image"
    "github.com/go-gl/mathgl/mgl32"
    "github.com/go-gl/gl/v4.5-core/gl"
    "strings"
    "image"
    "image/draw"
    "fmt"
    "bytes"
    "os"
    "io"
)

const (
)

type ImageTexture struct {
    TextureHandle uint32
    TextureSize   mgl32.Vec2
}



func compileShader(source string, shaderType uint32) (uint32, error) {
    shader := gl.CreateShader(shaderType)

    csources, free := gl.Strs(source)
    var csourceslength int32 = int32(len(source))
    gl.ShaderSource(shader, 1, csources, &csourceslength)
    free()
    gl.CompileShader(shader)

    var status int32
    gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

        return 0, fmt.Errorf("failed to compile %v: %v", source, log)
    }

    return shader, nil
}

func readFile(name string) (string, error) {

    buf := bytes.NewBuffer(nil)
    f, err := os.Open(name)
    if err != nil {
        return "", err
    }
    io.Copy(buf, f)
    f.Close()

    return string(buf.Bytes()), nil
}

// Mostly taken from the Demo. But compiling and linking shaders
// just should be done like this anyways.
func NewProgram(vertexShaderName, fragmentShaderName string) (uint32, error) {

    vertexShaderSource, err := readFile(vertexShaderName)
    if err != nil {
        return 0, err
    }

    fragmentShaderSource, err := readFile(fragmentShaderName)
    if err != nil {
        return 0, err
    }

    vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
    if err != nil {
        return 0, err
    }

    fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
    if err != nil {
        return 0, err
    }

    program := gl.CreateProgram()

    gl.AttachShader(program, vertexShader)
    gl.AttachShader(program, fragmentShader)
    gl.LinkProgram(program)

    var status int32
    gl.GetProgramiv(program, gl.LINK_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

        return 0, fmt.Errorf("failed to link program: %v", log)
    }

    gl.DeleteShader(vertexShader)
    gl.DeleteShader(fragmentShader)

    return program, nil
}

func NewComputeProgram(computeShaderName string) (uint32, error) {

    computeShaderSource, err := readFile(computeShaderName)
    if err != nil {
        return 0, err
    }
    computeShader, err := compileShader(computeShaderSource, gl.COMPUTE_SHADER)
    if err != nil {
        return 0, err
    }
    program := gl.CreateProgram()

    gl.AttachShader(program, computeShader)
    gl.LinkProgram(program)

    var status int32
    gl.GetProgramiv(program, gl.LINK_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

        return 0, fmt.Errorf("failed to link program: %v", log)
    }

    gl.DeleteShader(computeShader)

    return program, nil
}

func CreateTexture(width, height int32, internalFormat, format, internalType uint32, multisampling bool, samples, mipmapLevels int32) uint32 {

    var texType uint32 = gl.TEXTURE_2D
    if multisampling {
        texType = gl.TEXTURE_2D_MULTISAMPLE
    }

    var tex uint32
    gl.GenTextures(1, &tex)
    gl.BindTexture(texType, tex)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

    if multisampling {
        gl.TexStorage2DMultisample(gl.TEXTURE_2D_MULTISAMPLE, samples, internalFormat, width, height, false)
    } else {
        gl.TexStorage2D(gl.TEXTURE_2D, mipmapLevels, internalFormat, width, height)
    }

    return tex
}

func CreateImageTexture(imageName string, isRepeating bool) ImageTexture {

    var imageTexture ImageTexture

    img, err := LoadImage(imageName)
    if err != nil {
        fmt.Printf("Image load failed: %v.\n", err)
    }

    var textureWrap int32 = gl.CLAMP_TO_EDGE
    if isRepeating {
        textureWrap = gl.REPEAT
    }

    rgbaImg := image.NewRGBA(img.Img.Bounds())
    draw.Draw(rgbaImg, rgbaImg.Bounds(), img.Img, image.Pt(0, 0), draw.Src)

    gl.GenTextures(1, &imageTexture.TextureHandle)
    gl.BindTexture(gl.TEXTURE_2D, imageTexture.TextureHandle)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, textureWrap)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, textureWrap)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
    gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, int32(img.Img.Bounds().Max.X), int32(img.Img.Bounds().Max.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgbaImg.Pix))

    imageTexture.TextureSize = mgl32.Vec2{float32(img.Img.Bounds().Max.X), float32(img.Img.Bounds().Max.Y)}

    return imageTexture

}

func CreateFboWithExistingTextures(colorTex, depthTex *uint32, texType uint32) uint32{

    var fbo uint32
    gl.GenFramebuffers(1, &fbo)
    gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)

    if colorTex != nil {
        gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, texType, *colorTex, 0)
    }
    if depthTex != nil {
        gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT,  texType, *depthTex, 0)
    }

    gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

    return fbo
}

// Some internal format changes, like only having the RG channels but with higher 32F precision.
func CreateLightFbo(colorTex, depthTex *uint32, width, height int32, multisampling bool, samples int32) uint32 {

    if colorTex != nil {
        *colorTex = CreateTexture(width, height, gl.RG32F, gl.RG, gl.FLOAT, multisampling, samples, 1)
    }
    if depthTex != nil {
        *depthTex = CreateTexture(width, height, gl.DEPTH_COMPONENT32, gl.DEPTH_COMPONENT, gl.FLOAT, multisampling, samples, 1)
    }

    var texType uint32 = gl.TEXTURE_2D
    if multisampling {
        texType = gl.TEXTURE_2D_MULTISAMPLE
    }

    return CreateFboWithExistingTextures(colorTex, depthTex, texType)
}

func CreateFbo(colorTex, depthTex *uint32, width, height int32, multisampling bool, samples int32, isFloatingPoint bool, mipmapLevels int32) uint32 {

    var intFormat uint32 = uint32(gl.RGBA8)
    var format    uint32 = uint32(gl.RGBA)
    var ttype     uint32 = uint32(gl.UNSIGNED_BYTE)

    if isFloatingPoint {
        intFormat = gl.RGBA32F
        ttype = gl.FLOAT
    }

    if colorTex != nil {
        *colorTex = CreateTexture(width, height, intFormat, format, ttype, multisampling, samples, mipmapLevels)
    }
    if depthTex != nil {
        *depthTex = CreateTexture(width, height, gl.DEPTH_COMPONENT32, gl.DEPTH_COMPONENT, gl.FLOAT, multisampling, samples, 1)
    }

    var texType uint32 = gl.TEXTURE_2D
    if multisampling {
        texType = gl.TEXTURE_2D_MULTISAMPLE
    }

    return CreateFboWithExistingTextures(colorTex, depthTex, texType)
}


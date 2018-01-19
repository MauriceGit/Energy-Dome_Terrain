package opengl

import (
    //"github.com/go-gl/mathgl/mgl32"
    "github.com/go-gl/gl/v4.5-core/gl"
    "strings"
    "fmt"
    "bytes"
    "os"
    "io"
)

const (
)





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

func CreateMSTexture(tex *uint32, width, height, internalFormat int32, format, internalType uint32) {
    gl.GenTextures(1, tex);
    gl.BindTexture(gl.TEXTURE_2D_MULTISAMPLE, *tex);
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST);
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST);
    gl.TexImage2DMultisample(gl.TEXTURE_2D_MULTISAMPLE, 2, uint32(internalFormat), width, height, false);
}

func CreateTexture(tex *uint32, width, height, internalFormat int32, format, internalType uint32) {
    gl.GenTextures(1, tex);
    gl.BindTexture(gl.TEXTURE_2D, *tex);
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST);
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST);
    gl.TexImage2D(gl.TEXTURE_2D, 0, internalFormat, width, height, 0, format, internalType, nil);
}

func CreateFboWithExistingTextures(fbo, colorTex, depthTex *uint32, texType uint32) {
    gl.GenFramebuffers(1, fbo);
    gl.BindFramebuffer(gl.FRAMEBUFFER, *fbo);

    if colorTex != nil {
        gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, texType, *colorTex, 0);
    }
    if depthTex != nil {
        gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT,  texType, *depthTex, 0);
    }

    gl.BindFramebuffer(gl.FRAMEBUFFER, 0);
}


func CreateLightFbo(fbo, colorTex, depthTex *uint32, width, height int32, multisampling bool) {

    if colorTex != nil {
        if multisampling {
            CreateMSTexture(colorTex, width, height, gl.RG32F, gl.RG, gl.FLOAT)
        } else {
            CreateTexture(colorTex, width, height, gl.RG32F, gl.RG, gl.FLOAT)
        }
    }
    if depthTex != nil {
        if multisampling {
            CreateMSTexture(depthTex, width, height, gl.DEPTH_COMPONENT32, gl.DEPTH_COMPONENT, gl.FLOAT)
        } else {
            CreateTexture(depthTex, width, height, gl.DEPTH_COMPONENT32, gl.DEPTH_COMPONENT, gl.FLOAT)
        }
    }

    var texType uint32 = gl.TEXTURE_2D
    if multisampling {
        texType = gl.TEXTURE_2D_MULTISAMPLE
    }

    CreateFboWithExistingTextures(fbo, colorTex, depthTex, texType)

}

func CreateFbo(fbo, colorTex, depthTex *uint32, width, height int32, multisampling bool) {

    if colorTex != nil {
        if multisampling {
            CreateMSTexture(colorTex, width, height, gl.RGBA8, gl.RGBA, gl.UNSIGNED_BYTE)
        } else {
            CreateTexture(colorTex, width, height, gl.RGBA8, gl.RGBA, gl.UNSIGNED_BYTE)
        }
    }
    if depthTex != nil {
        if multisampling {
            CreateMSTexture(depthTex, width, height, gl.DEPTH_COMPONENT32, gl.DEPTH_COMPONENT, gl.FLOAT)
        } else {
            CreateTexture(depthTex, width, height, gl.DEPTH_COMPONENT32, gl.DEPTH_COMPONENT, gl.FLOAT)
        }
    }

    var texType uint32 = gl.TEXTURE_2D
    if multisampling {
        texType = gl.TEXTURE_2D_MULTISAMPLE
    }

    CreateFboWithExistingTextures(fbo, colorTex, depthTex, texType)

}





#!/bin/bash

export GOPATH=$(pwd)/Go
echo "Set the GOPATH to '$GOPATH'"
export GOBIN=$(pwd)/bin
echo "Set the GOBIN to '$GOBIN'"

if [[ ! -d "bin" ]]
then
    mkdir bin
else
    echo "./bin directory already exists"
fi

echo "Get Libaries"
go get "github.com/go-gl/gl/v4.5-core/gl"
go get "github.com/go-gl/glfw/v3.2/glfw"
go get "github.com/go-gl/mathgl/mgl32"

echo "Build Task"
go install HeightmapTerrain




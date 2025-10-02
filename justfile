
# Variables
binary_name := "data-octogon"
build_dir := "build"
version := "0.1.0"
commit := `git rev-parse --short HEAD`
date := `date -u +%Y-%m-%dT%H:%M:%SZ`
ldflags := "-X main.version={{version}} -X main.commit={{commit}} -X main.date={{date}}"
os := ["linux", "darwin"]
arch := ["amd64", "arm64"]

# Help function
default:
    just --list

# Installs all dependencies for the development container
install:
    #!/bin/bash

    bash init/post-devcontainer-create.sh

# Builds the data structures binary
build:
    mkdir -p {{build_dir}}
    for os in {{os}}; do
      for arch in {{arch}}; do
        echo "Building for $os/$arch..."
        GOOS=$os GOARCH=$arch go build -ldflags="{{ldflags}}" -o {{build_dir}}/{{binary_name}}-{{os}}-{{arch}}
      done
    done
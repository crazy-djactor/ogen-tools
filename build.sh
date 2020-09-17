#!/bin/bash

export COMPILER_NAME=ogen_compiler
export LAUNCHER_NAME=ogen_launcher

cd compiler && go build -o "$COMPILER_NAME" main.go && mv "$COMPILER_NAME" ../ && cd ..
cd launcher && go build -o "$ LAUNCHER_NAME" main.go && mv "$LAUNCHER_NAME" ../ && cd ..
#!/bin/bash

export COMPILER_NAME=ogen_compiler
export LAUNCHER_NAME=ogen_launcher
export COMPILER_FOLDER=ogen-compiler-linux_amd64
export LAUNCHER_FOLDER=ogen-launcher-linux_amd64

cd compiler && go build -o "$COMPILER_NAME" main.go && mv "$COMPILER_NAME" ../ && cd ..

rm -rf "$COMPILER_FOLDER"
mkdir "$COMPILER_FOLDER" && mv ./"$COMPILER_NAME" "$COMPILER_FOLDER"
tar -czvf "$COMPILER_FOLDER".tar.gz "$COMPILER_FOLDER"
rm -rf "$COMPILER_FOLDER"

cd launcher && go build -o "$LAUNCHER_NAME" main.go && mv "$LAUNCHER_NAME" ../ && cd ..

rm -rf "$LAUNCHER_FOLDER"
mkdir "$LAUNCHER_FOLDER" && mv ./"$LAUNCHER_NAME" "$LAUNCHER_FOLDER"
tar -czvf "$LAUNCHER_FOLDER".tar.gz "$LAUNCHER_FOLDER"
rm -rf "$LAUNCHER_FOLDER"
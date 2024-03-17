#!/bin/bash
ARR_OS=("linux")
ARR_ARCH=("arm" "amd64")

BUILD_PATH="./dist"
GO=$(which go)

export CGO_ENABLED="0"
export GOARM="6"
export GOPATH="/home/vizn3r/go"

/usr/bin/rm -rf "$BUILD_PATH/"

for ARCH in "${ARR_ARCH[@]}"
do
  for OS in "${ARR_OS[@]}"
  do
    FIRMWARE="eve-firmware-$OS-$ARCH"
    if [ "$OS" = "windows" ]
    then
      echo "BUILDING $BUILD_PATH/$FIRMWARE.exe"
      GOOS="$OS" GOARCH="$ARCH" /usr/bin/go/bin/go build -o "$BUILDPATH/$FIRMWARE.exe"
    else
      echo "BUILDING $PATH/$FIRMWARE"
      GOOS="$OS" GOARCH="$ARCH" /usr/bin/go/bin/go build -o "$BUILDPATH/$FIRMWARE"
    fi
  done
done

/usr/bin/pscp -pw 3766 "$BUILDPATH/eve-firmware-linux-arm" "simon@eve.local:"
# /usr/bin/pscp -pw 3766 "$PATH/eve-firmware-linux-arm64" "simon@eve.local:"

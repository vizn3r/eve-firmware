#!/bin/bash
ARR_OS=("linux")
ARR_ARCH=("arm" "amd64")

PATH="./dist"

export CGO_ENABLED="0"
export GOARM="6"

/usr/bin/rm -rf "$PATH/"

for ARCH in "${ARR_ARCH[@]}"
do
  for OS in "${ARR_OS[@]}"
  do
    FIRMWARE="eve-firmware-$OS-$ARCH"
    if [ "$OS" = "windows" ]
    then
      echo "BUILDING $PATH/$FIRMWARE.exe"
      GOOS="$OS" GOARCH="$ARCH" /usr/bin/go build -o "$PATH/$FIRMWARE.exe"
    else
      echo "BUILDING $PATH/$FIRMWARE"
      GOOS="$OS" GOARCH="$ARCH" /usr/bin/go build -o "$PATH/$FIRMWARE"
    fi
  done
done

/usr/bin/pscp -pw 3766 "$PATH/eve-firmware-linux-arm" "simon@eve.local:"
# /usr/bin/pscp -pw 3766 "$PATH/eve-firmware-linux-arm64" "simon@eve.local:"

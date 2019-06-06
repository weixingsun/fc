#!/bin/bash
compile() {
  ARCH=$1
  echo "Compiling for $ARCH"
  if [ "$ARCH" = "arm" ]; then
    echo "#raspberry pi zero"
    env GOOS=linux GOARCH=$ARCH GOARM=6 go build -o bin/fc_$ARCH
  else
    env GOOS=linux GOARCH=$ARCH go build -o bin/fc_$ARCH
  fi
}
find_arch() {
  arch=""
  case $(uname -m) in
    i386)   arch="386" ;;
    i686)   arch="386" ;;
    x86_64) arch="amd64" ;;
    arm)    dpkg --print-architecture | grep -q "arm64" && arch="arm64" || arch="arm" ;;
  esac
  echo $arch
  #return $arch
}
#find_arch 
#local_arch=$?
#compile $local_arch
compile arm64
compile arm
compile amd64

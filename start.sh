#!/bin/bash
arch=""
case $(uname -m) in
  i386)   arch="386" ;;
  i686)   arch="386" ;;
  x86_64) arch="amd64" ;;
  arm)    dpkg --print-architecture | grep -q "arm64" && arch="arm64" || arch="arm" ;;
esac
#echo $arch

SCRIPT=$(readlink -f $0)
DIR=`dirname $SCRIPT`
mkdir -p $DIR/log
nohup $DIR/bin/fc_$arch > $DIR/log/start.log 2>&1 &

#!/bin/bash

pid=`ps -ef|grep fc_go|grep -v grep|grep -v start.sh|awk '{print $2}'`
sudo kill -9 $pid
arch=""
case $(uname -m) in
  i386)   arch="386" ;;
  i686)   arch="386" ;;
  x86_64) arch="amd64" ;;
  armv6l) arch="arm" ;;
  arm)    dpkg --print-architecture | grep -q "arm64" && arch="arm64" || arch="arm" ;;
esac
#echo $arch

SCRIPT=$(readlink -f $0)
DIR=`dirname $SCRIPT`
mkdir -p $DIR/log
nohup $DIR/bin/fc_$arch > $DIR/log/start.log 2>&1 &

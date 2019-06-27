#!/bin/bash

set -e

runmode_from_conf=`awk '$1=="runmode" {print $3}' ./conf/app.conf`
version_from_conf=`awk '$1=="version" {print $3}' ./conf/app.conf`
apiport_from_conf=`awk '$1=="apiport" {print $3}' ./conf/app.conf`
socketport_from_conf=`awk '$1=="socketport" {print $3}' ./conf/app.conf`

if [ $# == 0 ] && [ -z $version_from_conf ]; then
    echo "baby, we need a version code"
    exit 1
fi

runmode=$runmode_from_conf
if [ $# == 1 ]; then
    runmode=$1
fi

echo $runmode

version=$version_from_conf
if [ $# == 1 ]; then
    version=$1
fi

echo $version

apiport=$apiport_from_conf
if [ $# == 1 ]; then
    apiport=$1
fi

echo $apiport

socketport=$socketport_from_conf
if [ $# == 1 ]; then
    socketport=$1
fi

echo $socketport

default_runmode="dev"
runmode=`awk '$1=="runmode" {print $3}' ./conf/app.conf`

if [ $default_runmode != $runmode ]
then
    echo "$runmode is err,you should in $default_runmode"
	exit 1
fi

ssh  root@106.14.202.179 version=$version apiport=$apiport socketport=$socketport runmode=$runmode 'bash -se' <<'ENDSSH'
cd ~/app/api/shici/dev/shici
git pull;
echo shici\_$runmode
#go clean;
if docker build -t shici\_$runmode:$version .
then
    echo "stop and rm old container,start new one..."
    docker stop shici\_$runmode
    docker rm shici\_$runmode
    docker run --restart=always --name shici\_$runmode -d -p $apiport:8080 -p $socketport:6666 shici\_$runmode:$version
    docker ps
fi
ENDSSH
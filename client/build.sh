#!/bin/bash
gobuild(){
	export GOOS=$1
	export GOARCH=$2
	export GOARM=$3
	sufix=""
	if [ "$GOOS" == "windows" ];then
		sufix=".exe"
	fi
	echo "build $GOOS-$GOARCH$GOARM$sufix"
	if [ "$GOARCH" == "arm" ];then
		go build -o "sow-client-$GOOS-$GOARCH$GOARM$sufix"
	else
		go build -o "sow-client-$GOOS-$GOARCH$sufix"
	fi
	move $4
}

move(){
	if [ -d "$1" ];then
		mv "client-$GOOS-$GOARCH" "$1"
	else
		echo "dir [$1] not exists"
	fi
}

build(){
	archs=("amd64" "386")
	gooses=("linux" "darwin" "windows")
	arms=(6 5)
	
	for goos in ${gooses[@]};do
		for arch in ${archs[@]};do
			gobuild $goos $arch "" $1
		done
	done
	gobuild "linux" "arm" "6" $1
	gobuild "linux" "arm" "5" $1
}

clean(){
	rm -rf sow-client-*
}

if [ "$1" == "clean" ];then
	clean
else
	build $2
fi




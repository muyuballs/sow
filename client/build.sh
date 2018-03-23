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
	target="sow-client-$GOOS-$GOARCH$sufix"
	if [ "$GOARCH" == "arm" ];then
		target="sow-client-$GOOS-$GOARCH$GOARM$sufix"
	fi
	go build -o $target
	move $target $4
}

move(){
	if [ -d "$2" ];then
		mv "$1" "$2"
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




PROTO_ROOT?=torntools_proto/torntools
PROTO_BUILD_DIR?=gen

include torntools_proto/Makefile

torntools_server:
	echo "build main server"
	mkdir -p build/bin
	go mod vendor
	go build  -o  build/bin/torntools_server ./server/cmd/main.go

simple_client:
	echo "build test client"
	mkdir -p build/bin
	go mod vendor
	go build  -o  build/bin/simple_client ./client/cmd/main.go

clean:
	rm -rf gen
	rm -rf build/bin

all: proto torntools_server
	echo "build all"
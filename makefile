.PHONY:  build-api build-web    build-cli help

apiBinName="ginskeleton-api.linux64"

webBinName="ginskeleton-web.linux64"

cliBinName="ginskeleton-cli.linux64"

all:
	go env -w GOARCH=amd64
	go env -w GOOS=linux
	go env -w CGO_ENABLED=0
	go env -w GO111MODULE=on
	go env -w GOPROXY=
	go mod  tidy

build-api:all clean-api build-api-bin
build-api-bin:
	go build -o ${apiBinName}    -ldflags "-w -s"  -trimpath  ./cmd/api/main.go

build-web:all clean-web build-web-bin
build-web-bin:
	go build -o ${webBinName}   -ldflags "-w -s"  -trimpath  ./cmd/web/main.go

build-cli:all clean-cli build-cli-bin
build-cli-bin:
	go build -o ${cliBinName}   -ldflags "-w -s"  -trimpath  ./cmd/cli/main.go

clean-api:
	@if [ -f ${apiBinName} ] ; then rm -rf ${apiBinName} ; fi
clean-web:
	@if [ -f ${webBinName} ] ; then rm -rf ${webBinName} ; fi
clean-cli:
	@if [ -f ${cliBinName} ] ; then rm -rf ${cliBinName} ; fi

help:
	@echo "make hep"
	@echo "make build-api cmd/api/main.go"
	@echo "make build-web cmd/web/main.go"
	@echo "make build-cli cmd/cli/main.go"
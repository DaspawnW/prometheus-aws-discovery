PROJECT_NAME := prometheus-aws-discovery
GITHUB_PATH := github.com/daspawnw/prometheus-aws-discovery

all: clean test bin

clean:
	rm -rf bin

test:
	go test ${GITHUB_PATH}/...

bin: bin/linux bin/darwin

bin/%:
	mkdir -p $@
	CGO_ENABLED=0 GOOS=$(word 1, $(subst /, ,$*)) GOARCH=amd64 go build -o "$@" ${GITHUB_PATH}/cmd/${PROJECT_NAME}/...

run:
	go run ${GITHUB_PATH}/cmd/${PROJECT_NAME}/...
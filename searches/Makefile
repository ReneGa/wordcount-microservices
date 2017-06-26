.PHONY: clean

all: searches

clean:
	rm -f searches

BUILD_DOCKER_IMAGE_NAME=wordcount-microservices-searches-build
GO_PACKAGE_PATH=/go/src/github.com/ReneGa/tweetcount-microservices/searches
BUILD_COMMAND=go get && go build  -ldflags '-linkmode external -extldflags -static'

searches: $(shell find . -iname '*.go')
	docker build -t ${BUILD_DOCKER_IMAGE_NAME} -f Dockerfile_build .
	docker run --rm -v "${PWD}":"${GO_PACKAGE_PATH}" -w "${GO_PACKAGE_PATH}" ${BUILD_DOCKER_IMAGE_NAME} /bin/sh -c "${BUILD_COMMAND}"
.PHONY: clean

all: searches

clean:
	rm -f searches

BUILD_DOCKER_IMAGE_NAME=wordcount-microservices-searches-build
PACKAGE_PATH_IN_IMAGE=/go/src/github.com/ReneGa/tweetcount-microservices/searches
GO_BUILD_COMMAND=go get && go build  -ldflags '-linkmode external -extldflags -static'

searches: $(shell find . -iname '*.go')
	docker build -t ${BUILD_DOCKER_IMAGE_NAME} -f Dockerfile_build .
	docker run --rm -v "${PWD}":"${PACKAGE_PATH_IN_IMAGE}" -w "${PACKAGE_PATH_IN_IMAGE}" ${BUILD_DOCKER_IMAGE_NAME} /bin/sh -c "${GO_BUILD_COMMAND}"
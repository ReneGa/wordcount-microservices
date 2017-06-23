.PHONY: all clean

all: ingestor

clean:
	rm -f ingestor

ingestor: $(shell find . -iname '*.go')
	GOOS=linux go build
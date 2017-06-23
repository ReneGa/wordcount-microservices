.PHONY: all clean

all: wordcounter

clean:
	rm -f wordcounter

wordcounter: $(shell find . -iname '*.go')
	GOOS=linux go build
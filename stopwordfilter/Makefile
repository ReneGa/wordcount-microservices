.PHONY: clean

all: stopwordfilter

clean:
	rm -f stopwordfilter

stopwordfilter: $(shell find . -iname '*.go')
	GOOS=linux go build
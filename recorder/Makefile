.PHONY: clean

all: recorder

clean:
	rm -f recorder

recorder: $(shell find . -iname '*.go')
	GOOS=linux go build
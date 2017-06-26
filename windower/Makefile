.PHONY: clean

all: windower

clean:
	rm -f windower

windower: $(shell find . -iname '*.go')
	GOOS=linux go build
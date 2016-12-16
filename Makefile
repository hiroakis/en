.PHONY: all clean build

all: clean build

test:
	go test -v

build: test
	go build

clean:
	rm -f en


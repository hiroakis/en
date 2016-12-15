
test:
	go test -v

build: test
	go build

clean:
	rm -f en

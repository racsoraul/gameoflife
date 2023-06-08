run: clean build
	bin/gameoflife

build:
	go build -o bin/gameoflife ./

clean:
	rm -f bin/*

PHONY: build clean
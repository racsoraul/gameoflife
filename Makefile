run_dev: clean build_dev
	bin/gameoflife_dev

build_dev:
	go build -race -o bin/gameoflife_dev ./

run: clean build
	bin/gameoflife

build:
	go build -o bin/gameoflife ./

clean:
	rm -f bin/*

PHONY: run_dev build_dev run clean build
UNAME_S := $(shell uname -s)

# go-sdl2's ttf/sdl_ttf_cgo.go specifies both pkg-config and explicit LDFLAGS
# for SDL2_ttf, causing duplicate -l flags. This is harmless but produces a
# warning on macOS. Suppress it only on macOS where the flag is supported.
ifeq ($(UNAME_S),Darwin)
CGO_LDFLAGS_EXTRA := -Wl,-no_warn_duplicate_libraries
endif

run: clean build
	bin/gameoflife
build:
	CGO_LDFLAGS="$(CGO_LDFLAGS_EXTRA)" go build -o bin/gameoflife ./
run_dev: clean build_dev
	bin/gameoflife_dev
build_dev:
	CGO_LDFLAGS="$(CGO_LDFLAGS_EXTRA)" go build -race -o bin/gameoflife_dev ./
clean:
	rm -f bin/*

.PHONY: run_dev build_dev run clean build

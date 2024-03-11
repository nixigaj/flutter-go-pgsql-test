.PHONY: copy-static build build-release run clean

ifeq ($(shell uname -s),Linux)
CGO_ENABLED := 0
endif

all: build

copy-static:
	@mkdir -p build
	@if [ ! -f "build/config.toml" ]; then \
		cp static/config.toml build/; \
	fi

build: copy-static
	@CGO_ENABLED=$(CGO_ENABLED) go build -o build/backend cmd/*.go

build-release: copy-static
	@CGO_ENABLED=$(CGO_ENABLED) go build --ldflags "-w -s" -o build/backend cmd/*.go

run:
	@cd build && ./backend && cd ..

clean:
	@rm -rf build/

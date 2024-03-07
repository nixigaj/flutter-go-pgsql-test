.PHONY: copy-static build build-release run clean

all: build

copy-static:
	@mkdir -p build
	@if [ ! -f "build/config.toml" ]; then \
		cp static/config.toml build/; \
	fi

build: copy-static
	@go build -o build/backend cmd/*.go

build-release: copy-static
	@go build --ldflags "-w -s" -o build/backend cmd/*.go

run:
	@cd build && ./backend && cd ..

clean:
	@rm -rf build/

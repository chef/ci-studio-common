LD_FLAGS=-s -w
UTILITIES = ci-studio-common-util file-mod install-buildkite-agent install-habitat did-modify vault-util

.PHONE: studio
studio:
	HAB_ORIGIN='chef' hab studio -k chef enter

.PHONY: build
build: build-darwin build-linux build-windows

.PHONY: clean-all clean
clean-all:
	rm -rf build

clean: clean-linux clean-windows clean-darwin

clean-%:
	rm -rf build/$*

#
# Darwin
#
.PHONY: build-darwin
build-darwin: $(addprefix build/darwin/,$(UTILITIES))

build/darwin/%:
	mkdir -p build/darwin
	GOOS=darwin GOARCH=amd64 go build -o build/darwin/$* -ldflags="$(LD_FLAGS)" ./cmd/$*

#
# Linux
#
.PHONY: build-linux
build-linux: $(addprefix build/linux/,$(UTILITIES))

build/linux/%:
	mkdir -p build/linux
	GOOS=linux GOARCH=amd64 go build -o build/linux/$* -ldflags="$(LD_FLAGS)" ./cmd/$*

#
# Windows
#
.PHONY: build-windows
build-windows: $(addsuffix .exe, $(addprefix build/windows/,$(UTILITIES)))

build/windows/%.exe:
	mkdir -p build/windows
	GOOS=windows GOARCH=amd64 go build -o build/windows/$*.exe -ldflags="$(LD_FLAGS)" ./cmd/$*

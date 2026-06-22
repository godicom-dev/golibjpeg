GOHOSTOS ?= $(shell go env GOHOSTOS)
GOHOSTARCH ?= $(shell go env GOHOSTARCH)

# Map Go arch → cmake arch
CMAKE_ARCH_amd64 = x64
CMAKE_ARCH_arm64 = arm64
CMAKE_ARCH_386   = Win32

embed-name: ## print the native lib filename for the current host
ifeq ($(GOHOSTOS),windows)
	@echo golibjpeg_$(GOHOSTARCH).dll
else ifeq ($(GOHOSTOS),darwin)
	@echo golibjpeg_darwin_$(GOHOSTARCH).dylib
else
	@echo golibjpeg_linux_$(GOHOSTARCH).so
endif

.PHONY: build-native
build-native: ## build native library for the current host
	cmake -S _csrc -B _csrc/build \
	  -DCMAKE_BUILD_TYPE=Release \
	  -DCMAKE_POSITION_INDEPENDENT_CODE=ON
	cmake --build _csrc/build --config Release
	mkdir -p native/libs
	cp _csrc/build/libgolibjpeg.* native/libs/

.PHONY: build-native-win
build-native-win: ## build native DLL on Windows with MSVC
	cmake -S _csrc -B _csrc/build \
	  -G "Visual Studio 17 2022" -A $(CMAKE_ARCH_$(GOHOSTARCH))
	cmake --build _csrc/build --config Release
	mkdir -p native/libs
	cp _csrc/build/Release/golibjpeg.dll native/libs/golibjpeg_$(GOHOSTARCH).dll

.PHONY: test
test: ## run Go tests
	go test -v ./...

.PHONY: clean
clean: ## remove build artifacts
	rm -rf _csrc/build _csrc/build_nix

.PHONY: tidy
tidy: ## tidy Go module
	go mod tidy

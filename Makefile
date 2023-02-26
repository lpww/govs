build:
	@echo "Building from source"
	go build

clean:
	@echo "Cleaning build artifacts"
	rm -rf build govs

install:
	@echo "Installing modules"
	go get

release:
	@echo "Creating Release Binaries"
	go install github.com/mitchellh/gox@latest
	mkdir -p build
	gox -output="build/{{.OS}}_{{.Arch}}/{{.Dir}}" -osarch="darwin/amd64 linux/amd64 windows/amd64 openbsd/amd64"
	tar czf build/darwin_amd64.tar.gz build/darwin_amd64
	tar czf build/linux_amd64.tar.gz build/linux_amd64
	tar czf build/windows_amd64.tar.gz build/windows_amd64
	tar czf build/openbsd_amd64.tar.gz build/openbsd_amd64

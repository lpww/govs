build:
	@echo "Building from source"
	go build

install:
	@echo "Installing modules"
	go get

release:
	@echo "Creating Release Binaries"
	go install github.com/mitchellh/gox@latest
	mkdir -p build
	gox -output="build/{{.OS}}_{{.Arch}}/{{.Dir}}"


all: setup_deps

setup_deps:
	go mod tidy
	go mod vendor
	#go build

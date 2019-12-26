# Go parameters
# SHELL := cmd.exe
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
SRC=./src/
BINARY_DIR=./bin/

BINARY_NAME=rein-1.0.6
BINARY_AMD64_WIN=$(BINARY_NAME)-amd64-win.exe
BINARY_AMD64_MAC=$(BINARY_NAME)-amd64-mac.dmg
BINARY_AMD64_LINUX=$(BINARY_NAME)-amd64-linux
BINARY_386_WIN=$(BINARY_NAME)-386-win.exe
BINARY_386_LINUX=$(BINARY_NAME)-386-linux
BINARY_ARM_LINUX=$(BINARY_NAME)-arm-linux

all: clean build
all-platform: clean build-linux build-linux-386 build-windows build-windows-386 build-mac build-arm
test: clean build run
build:
		$(GOBUILD) -o $(BINARY_DIR)$(BINARY_NAME) -v $(SRC)
clean:
		rm -rf $(BINARY_DIR)*
run1:
		#$(GOBUILD) -o $(BINARY_DIR)$(BINARY_NAME) -v $(SRC)
run:
		$(BINARY_DIR)$(BINARY_NAME)
deps:
		#$(GOGET) github.com/markbates/goth
		#$(GOGET) github.com/markbates/pop
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)$(BINARY_AMD64_LINUX) -v $(SRC)
build-linux-386:
		CGO_ENABLED=0 GOOS=linux GOARCH=386 $(GOBUILD) -o $(BINARY_DIR)$(BINARY_386_LINUX) -v $(SRC)
build-windows:
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)$(BINARY_AMD64_WIN) -v $(SRC)
build-windows-386:
		CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GOBUILD) -o $(BINARY_DIR)$(BINARY_386_WIN) -v $(SRC)
build-mac:
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)$(BINARY_AMD64_MAC) -v $(SRC)
build-arm:
		CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GOBUILD) -o $(BINARY_DIR)$(BINARY_ARM_LINUX) -v $(SRC)


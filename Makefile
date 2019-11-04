# Go parameters
# SHELL := cmd.exe
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
SRC=./src/
BINARY_DIR=./bin/

BINARY_NAME=rein-1.0.5
BINARY_NAME_POSTFIX=.exe
BINARY_WIN=$(BINARY_NAME)-amd64-win.exe
BINARY_MAC=$(BINARY_NAME)-amd64-mac.dmg
BINARY_UNIX=$(BINARY_NAME)-amd64-linux


all: clean build
all-platform: clean build-linux build-windows build-mac
test: clean build run
build:
		$(GOBUILD) -o $(BINARY_DIR)$(BINARY_NAME)$(BINARY_NAME_POSTFIX) -v $(SRC)
clean:
		rm -rf $(BINARY_DIR)*
run1:
		#$(GOBUILD) -o $(BINARY_DIR)$(BINARY_NAME)$(BINARY_NAME_POSTFIX) -v $(SRC)
run:
		$(BINARY_DIR)$(BINARY_NAME)$(BINARY_NAME_POSTFIX)
deps:
		#$(GOGET) github.com/markbates/goth
		#$(GOGET) github.com/markbates/pop
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)$(BINARY_UNIX) -v $(SRC)
build-windows:
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)$(BINARY_WIN) -v $(SRC)
build-mac:
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)$(BINARY_MAC) -v $(SRC)


# import environment variables from ".env" file;
# ".env" file not stored in git, but ".env.example" is stored in git for reference; 
include ./.env

all: run

clean:
	rm -rf dist/

run: 
# run server from dist directory to easy access public files and to resemble situation after deployment;
# for performance reasons, we don't generate css at development time for every run;
	# run ##################################################
	mkdir -p ./dist
	cp -r ./src/resources ./dist
	(cd ./dist && go run ../src) 

build: clean 
# create dist directory from scratch with linux executable, resources and public files;
# resulting directory can be used for linux oci deployment; 
	# build ##################################################
	cp -r ./src/resources ./dist
	go version
# GOOS=windows GOARCH=amd64 go build -o dist/$(OCI_NAME)-windows.exe ./src
# file dist/$(OCI_NAME)-windows.exe
# GOOS=linux GOARCH=amd64 go build -o dist/$(OCI_NAME)-linux.exe ./src
# file dist/$(OCI_NAME)-linux.exe
	GOOS=linux GOARCH=amd64 go build -o dist/main.exe ./src
# Cross-Compilation:
# Install gcc-mingw toolchain on Ubuntu 24 for Windows cross-compilation:
#	sudo apt-get update && sudo apt-get install gcc-mingw-w64
# Install tdm-gcc toolchain on Windows 11 with WSL2 for Linux cross-compilation:
#   # Download and install tdm-gcc from https://jmeubank.github.io/tdm-gcc/
# Run command: (cd dist && ./main.exe)


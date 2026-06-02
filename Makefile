# import environment variables from ".env" file;
# ".env" file not stored in git, but ".env.example" is stored in git for reference
include ./.env

# makefile is used for automating tasks, not for file dependencies
.PHONY: $(wildcard *)

all: run

# --------------------------------------------------------------------------------

## help: Show this help message
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' Makefile | column -t -s ':' | sed -e 's/^/  /'
	@echo ""
	@echo "Example workflows:"
	@echo "> development support:    db-run ... run & db-migrate (*) ... db-stop" 
	@echo "> deployment preparation: oci-build ... db-run & db-migrate ... oci-run ... db-stop ... db-build"

# --------------------------------------------------------------------------------

## clean: Clean the dist directory
clean:
	rm -rf dist/
	mkdir -p ./dist

# --------------------------------------------------------------------------------

## build: Build the application executable to be used for Linux deployment  
## : (no .env configuration is used)
build: clean 
# alternative for watched generation of css using "npm run dev"
	npx @tailwindcss/cli -i ./src/resources/style.tailwindcss.css -o ./dist/public/style.css
# build linux executable in dist directory
	cp -r ./src/resources ./dist
	go mod tidy
	go version
# GOOS=windows GOARCH=amd64 go build -o dist/$(APP_NAME)-windows.exe ./src
# file dist/$(APP_NAME)-windows.exe
# GOOS=linux GOARCH=amd64 go build -o dist/$(APP_NAME)-linux.exe ./src
# file dist/$(APP_NAME)-linux.exe
	GOOS=linux GOARCH=amd64 go build -o dist/main.exe ./src
# Cross-Compilation:
# Install gcc-mingw toolchain on Ubuntu 24 for Windows cross-compilation:
#	sudo apt-get update && sudo apt-get install gcc-mingw-w64
# Install tdm-gcc toolchain on Windows 11 with WSL2 for Linux cross-compilation:
#   # Download and install tdm-gcc from https://jmeubank.github.io/tdm-gcc/
# Run command: (cd dist && ./main.exe)

# --------------------------------------------------------------------------------

## run: Run the application from source code for development purposes
## : (database must be accessible and .env configuration is used)
run: clean
# alternative for watched generation of css using "npm run dev"
	npx @tailwindcss/cli -i ./src/resources/style.tailwindcss.css -o ./dist/public/style.css
# run server from dist directory to easy access public files and to resemble situation after deployment
	cp -r ./src/resources ./dist
	cp ./.env ./dist/.env
	(cd ./dist && go run ../src) 
	rm ./dist/.env

# --------------------------------------------------------------------------------

# current time for labeling docker images with build time;
MAKE_TIME = $(shell date +"%FT%H:%M:%SZ")

## oci-build: Build OCI image using Docker and save as tar file
oci-build: build
	docker build . -t $(APP_NAME):latest --label "version=${APP_VERSION}" --label "build=$(MAKE_TIME)" \
	  --label "maintainer=https://github.com/RobertTC32/app-$(APP_NAME)"
	docker image ls | grep $(APP_NAME) 
	docker image inspect $(APP_NAME):latest
# docker rmi $(APP_NAME):latest
	mkdir -p ./dist
	docker save -o dist/$(APP_NAME).tar $(APP_NAME):latest
	du -sh dist/$(APP_NAME).tar
# docker load -i dist/$(APP_NAME).tar

# --------------------------------------------------------------------------------

## oci-run: Run existing OCI image using Docker
## : (database must be accessible and .env configuration is used)
oci-run: 
	docker run  --name myapp --rm -p $(APP_PORT):80 --env-file ./.env $(APP_NAME):latest

# --------------------------------------------------------------------------------

## db-clean: Clean the dist-db directory
db-clean:
	rm -rf dist-db/
	mkdir -p ./dist-db
	
# --------------------------------------------------------------------------------

## db-build: Build database resources & migration files to be used for deployment
db-build: db-clean
# create dist-db directory from scratch with database resources and migration files
	mkdir -p ./dist-db/db
	cp -r ./src-db/resources ./dist-db
	cp -r ./src-db/migrations ./dist-db
	cp ../goose.exe ./dist-db

# --------------------------------------------------------------------------------

## db-run: Create and run database for development purposes
## : (.env configuration is used)
db-run:  
# create test database in dist-db directory from scratch
	rm -rf temp/db
	mkdir -p ./temp/db
	# start mariadb container with mounted volumes for database and migration files
	docker run --name mariadb --rm -d -v ./temp/db:/var/lib/mysql:Z -p $(OCI_DB_PORT):3306 --env-file ./.env mariadb:12.2 
# docker run --name mariadb --rm -d -v ./temp/db:/var/lib/mysql:Z  -v ./src-db/resources:/docker-entrypoint-initdb.d \
#	  -p $(OCI_DB_PORT):3306 --env-file ./.env mariadb:12.2 
# docker exec -it mydb mariadb -uroot -p

# --------------------------------------------------------------------------------

## db-migrate: Migrate existing database for development purposes
## : (database must be accessible and .env configuration is used)
db-migrate:
# migrate database using goose sql scripts in src-db/migrations/ directory
# (create goose sql scripts using "goose -s create ${name} sql" command)
	goose up 

# --------------------------------------------------------------------------------

## db-stop: Stop running database for development purposes
## : (database must be running and no .env configuration is used)
db-stop:
	docker stop mariadb 

# --------------------------------------------------------------------------------


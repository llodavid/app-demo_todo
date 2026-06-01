# import environment variables from ".env" file;
# ".env" file not stored in git, but ".env.example" is stored in git for reference; 
include ./.env

MAKE_TIME = $(shell date +"%FT%H:%M:%SZ")

all: run

clean:
	rm -rf dist/

generate-twc: 
# alternative for watched generation of css using "npm run dev"; 
	npx @tailwindcss/cli -i ./src/resources/style.tailwindcss.css -o ./dist/public/style.css

run: 
# run server from dist directory to easy access public files and to resemble situation after deployment;
# for performance reasons, we don't generate css at development time for every run;
	# run ##################################################
	mkdir -p ./dist
	cp -r ./src/resources ./dist
	(cd ./dist && go run ../src) 

build: clean generate-twc 
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

clean-db:
	rm -rf dist-db/
	mkdir -p ./dist-db

test-db: 
# create test database in temp directory from scratch;
	# test-db ##################################################
	rm -rf ./temp/
	mkdir -p ./temp/db
	# sql scripts in src-db/resources/ are used for database migration
	cp -r ./src-db/resources ./temp
	# start mariadb container with mounted volumes for database and migration files;
	docker run --name mydb --rm -d -v ./temp/db:/var/lib/mysql:Z  -v ./temp/resources:/docker-entrypoint-initdb.d -p $(DB_PORT):3306 -e MARIADB_USER=$(DB_USER) -e MARIADB_PASSWORD=$(DB_PASSWORD) -e MARIADB_ROOT_PASSWORD=$(DB_ROOT_PASSWORD) -e MARIADB_DATABASE=$(DB_NAME) mariadb:12.2 
# docker exec -it mydb mariadb -uroot -p

build-db: clean-db
# create dist-db directory from scratch with database and migration files;
# resulting directory can be used for linux database deployment; 
	# build-db ##################################################
	mkdir -p ./dist-db/db
	cp -r ./src-db/resources ./dist-db

build-oci: build
	# build-oci ##################################################
	docker build . -t $(OCI_NAME):latest --label "version=${OCI_VERSION}" --label "build=$(MAKE_TIME)"
	docker image ls | grep $(OCI_NAME) 
	docker image inspect $(OCI_NAME):latest
	
run-oci: 
	# run-oci ##################################################
	docker run  --name myapp --rm -p $(OCI_PORT):$(OCI_INT_PORT) -e OCI_PORT=$(OCI_PORT) -e DB_PORT=$(DB_PORT) -e DB_HOST=host.docker.internal -e DB_USER=$(DB_USER) -e DB_PASSWORD=$(DB_PASSWORD) -e DB_NAME=$(DB_NAME) -e LOG_LEVEL=debug $(OCI_NAME):latest

save-oci:
	# save-oci ##################################################
	mkdir -p ./temp
	docker save -o temp/$(OCI_NAME).tar $(OCI_NAME):latest
	du -sh temp/$(OCI_NAME).tar
# docker rmi $(OCI_NAME):latest

load-oci:
	# load-oci ##################################################
	docker load -i temp/$(OCI_NAME).tar
	docker image ls | grep $(OCI_NAME)


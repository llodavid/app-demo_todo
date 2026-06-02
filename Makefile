# import environment variables from ".env" file;
# ".env" file not stored in git, but ".env.example" is stored in git for reference; 
include ./.env

MAKE_TIME = $(shell date +"%FT%H:%M:%SZ")

all: run

clean:
	# clean ##################################################
	rm -rf dist/
	mkdir -p ./dist

generate-twc: 
# alternative for watched generation of css using "npm run dev"; 
	npx @tailwindcss/cli -i ./src/resources/style.tailwindcss.css -o ./dist/public/style.css

run: generate-twc
# run server from dist directory to easy access public files and to resemble situation after deployment;
# for performance reasons, we don't generate css at development time for every run;
	# run ##################################################
	mkdir -p ./dist
	cp -r ./src/resources ./dist
	cp ./.env.example ./dist/.env
	(cd ./dist && go run ../src) 

build: clean generate-twc 
# create dist directory from scratch with linux executable, resources and public files;
# resulting directory can be used for linux oci deployment; 
	# build ##################################################
	cp -r ./src/resources ./dist
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

clean-db:
	# clean-db ##################################################
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
	docker run --name mydb --rm -d -v ./temp/db:/var/lib/mysql:Z  -v ./temp/resources:/docker-entrypoint-initdb.d -p $(OCI_DB_PORT):3306 --env-file ./.env mariadb:12.2 
# docker exec -it mydb mariadb -uroot -p

build-db: clean-db
# create dist-db directory from scratch with database and migration files;
# resulting directory can be used for linux database deployment; 
	# build-db ##################################################
	mkdir -p ./dist-db/db
	cp -r ./src-db/resources ./dist-db

build-oci: build
	# build-oci ##################################################
	docker build . -t $(APP_NAME):latest --label "version=${APP_VERSION}" --label "build=$(MAKE_TIME)"
	docker image ls | grep $(APP_NAME) 
	docker image inspect $(APP_NAME):latest
	
run-oci: 
	# run-oci ##################################################
	docker run  --name myapp --rm -p $(APP_PORT):80 -e APP_PORT=$(APP_PORT) -e DB_DSN="$(DB_DSN_OCI)" -e LOG_LEVEL=debug $(APP_NAME):latest

save-oci:
	# save-oci ##################################################
	mkdir -p ./temp
	docker save -o temp/$(APP_NAME).tar $(APP_NAME):latest
	du -sh temp/$(APP_NAME).tar
# docker rmi $(APP_NAME):latest

load-oci:
	# load-oci ##################################################
	docker load -i temp/$(APP_NAME).tar
	docker image ls | grep $(APP_NAME)


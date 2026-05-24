include ./.env

all: run

generate:
	npx @tailwindcss/cli -i ./src/templates/style-input.css -o ./public/css/style.css

run: generate
	go run ./src 

build: generate
	GOARCH=amd64 GOOS=windows go build -o bin/$(BINARY_NAME)-win.exe ./src
	GOARCH=amd64 GOOS=linux go build -o bin/$(BINARY_NAME)-lnx.exe ./src

clean:
	rm -rf bin/
	rm -rf public/

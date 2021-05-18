.PHONY: build clean deploy gomodgen remove dev

AWS_REGION ?= eu-west-3
SMARTHOME_JWT_EXPIRATION ?= 15m

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/setroom SetRoom/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getroom GetRoom/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/login Login/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/signup SignUp/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy -r $(AWS_REGION) --verbose
	go mod tidy

remove:
	sls remove

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

dev:
	go run main.go serve -a 127.0.0.1 --jwt-expiration $(SMARTHOME_JWT_EXPIRATION)

.PHONY: build-dev build clean deploy gomodgen remove dev

AWS_REGION ?= eu-west-3
SMARTHOME_JWT_SECRET ?= secret

build-dev:
	go mod tidy
	go build

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/setroom SetRoom/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getroom GetRoom/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/login Login/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/signup SignUp/main.go

clean:
	rm -rf ./bin ./vendor go.sum .serverless
	rm -f smarthome
	rm -f terraform/*.tfstate*
	rm -f terraform/.terraform.lock.hcl
	rm -rf terraform/.terraform
	docker compose down -v
	rm -rf docker
	go mod tidy

deploy: clean build
	sls deploy -r $(AWS_REGION) --verbose
	go mod tidy

remove:
	sls remove

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

dev: build-dev
	mkdir -p docker/dynamodb
	docker compose up -d
	terraform -chdir=terraform init
	terraform -chdir=terraform apply -auto-approve
	SMARTHOME_JWT_SECRET=$(SMARTHOME_JWT_SECRET) ./smarthome serve -a 127.0.0.1 \
		-v -d http://localhost:8000 \
		--dynamodb-auth-table Authentication \
		--dynamodb-control-table ControlPlane \
		--dynamodb-outside-table TemperatureOutside \
		--dynamodb-inside-table TemperatureInside
		--jwt-expiration 1h

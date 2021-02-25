.PHONY: build clean deploy gomodgen remove

AWS_REGION := eu-west-3

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/setroom SetRoom/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getroom GetRoom/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy -r $(AWS_REGION) --verbose

remove:
	sls remove

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

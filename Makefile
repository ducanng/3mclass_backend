dep:
	go mod tidy
	go mod download
build:
	go build -v -o ./bin/main .
build_deploy:
	go build -o functions/main ./...
run:
	./bin/main
generate:
	swag init --ot go,yaml -d ./
	swag fmt
clean:
	rm -rf ./functions
all:
	make clean
	make generate
	make dep
	make build
	make run
deploy:
	make clean
	go get ./...
	go install ./...
	make build_deploy
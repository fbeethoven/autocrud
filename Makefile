.PHONY: run clean test testv

all: run clean

run: 
	go run ./src/main.go ;
	go run ./myapitest/backend/src/main.go ;

clean: 
	rm -r myapitest

test: 
	go test ./... ;

testv: 
	go test -v ./... ;

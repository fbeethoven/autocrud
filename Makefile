.PHONY: run clean test testv

all: run clean

run: 
	go run ./src/main.go ;

exec:
	cd myapitest/backend/ && go run ./src/main.go ;

clean: 
	rm -r myapitest

test: 
	go test ./... ;

testv: 
	go test -v ./... ;

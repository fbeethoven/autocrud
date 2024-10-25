.PHONY: run clean test testv

all: run clean

run: 
	go run . -f config_test.yaml;

bexec:
	cd myapitest/backend/ && go run ./src/main.go ;

fexec:
	cd myapitest/frontend/ && npm run dev ;

clean: 
	rm -r myapitest

test: 
	go test ./... ;

testv: 
	go test -v ./... ;

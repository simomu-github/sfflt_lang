.PHONY: deps test clean

sfflt: deps
	go build -o release/sfflt cmd/sfflt.go

deps:
	go mod download

test: deps
	go test -v ./...

clean:
	rm -r release/*


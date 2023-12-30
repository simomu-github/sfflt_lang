.PHONY: deps test clean

sfflt: deps
	go build -o release/sfflt cmd/sfflt.go

run_test_script: sfflt
	./release/sfflt test.sflt
	$(FFLT_LANG) test.fflt

deps:
	go mod download

test: deps
	go test -v ./...

clean:
	rm -r release/*


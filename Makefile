bin/startc: cmd/startc/main.go pkg/namespaces/*
	go build -o bin/startc cmd/startc/main.go

.PHONY: build
build: bin/*

.PHONY: rsp
rsp:
	GOOS=linux GOARCH=arm64 make build

.PHONY: run
run: rsp
	rsync bin/startc rsproot:~/bin/
	ssh rsproot "~/bin/startc"

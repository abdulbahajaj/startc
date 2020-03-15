bin/startc: cmd/startc/main.go pkg/namespaces/* pkg/networking/*
	go build -o bin/startc cmd/startc/main.go

bin/netinit: cmd/netinit/main.go pkg/networking/*
	go build -o bin/netinit cmd/netinit/main.go
	chown root bin/netinit
	chmod u+s bin/netinit

.PHONY: build
build: bin/*

.PHONY: rsp
rsp:
	GOOS=linux GOARCH=arm64 make build

.PHONY: run
run: build
	bin/startc

.PHONY: runnet
runnet: build
	bin/netinit

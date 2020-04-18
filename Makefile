bin/startc: cmd/startc/main.go pkg/namespaces/* pkg/networking/* pkg/cgroup/*
	go build -o bin/startc cmd/startc/main.go

bin/pinit: cmd/privilagedInit/main.go pkg/networking/* pkg/cgroup/*
	go build -o bin/pinit cmd/privilagedInit/main.go
	sudo chown root bin/pinit
	sudo chmod u+s bin/pinit

.PHONY: build
build: bin/*

.PHONY: rsp
rsp:
	GOOS=linux GOARCH=arm64 make build

.PHONY: run
run: build
	bin/startc


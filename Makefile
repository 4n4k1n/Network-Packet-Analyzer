all:
	go build -o analyzer src/*.go

deps:
	go mod download github.com/google/gopacket

.PHONY: all deps
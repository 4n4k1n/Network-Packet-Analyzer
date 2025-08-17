all:
	go build -o analyzer src/*.go

run:
	go run src/*.c

deps:
	go mod download github.com/google/gopacket

.PHONY: all deps
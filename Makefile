all:
	go build -o analyzer src/*.go

run:
	go run src/*.c

update-mods:
	go mod tidy

deps:
	go mod download github.com/google/gopacket

.PHONY: all deps
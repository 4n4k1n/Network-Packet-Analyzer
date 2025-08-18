all:
	go build -o analyzer src/*.go

run:
	go run src/*.c

update-mods:
	go mod tidy

deps:
	go mod download github.com/google/gopacket
	go mod download github.com/dariubs/percent

clean:
	rm -f analyzer
	go clean
	go mod tidy

.PHONY: all deps clean
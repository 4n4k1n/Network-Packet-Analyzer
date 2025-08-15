FROM golang:latest

WORKDIR /app

RUN apt-get update && apt-get install -y libpcap-dev

COPY . .

RUN go mod download && go build -o network-analyzer .

EXPOSE 8080

CMD ["./network-analyzer"]
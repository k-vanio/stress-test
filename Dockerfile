FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o build/main ./cmd/main.go

RUN chmod +x ./build/main

CMD ["./build/main"]

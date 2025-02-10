FROM golang:1.22-alpine
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go test -v ./...
RUN go build -o main cmd/balance-service/main.go
CMD ["/app/main"]
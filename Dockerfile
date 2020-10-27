FROM golang:1.15-alpine

WORKDIR /go/src/
COPY . .
RUN GOOS=linux go build -ldflags="-s -w"
CMD ["./checkout-service"]
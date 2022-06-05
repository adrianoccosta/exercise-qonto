############################
# STEP 1 build executable binary
############################
FROM golang:1.18.3 AS builder

WORKDIR /go/src/qonto-service

COPY . .

RUN go get -u github.com/swaggo/swag/cmd/swag
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go mod vendor
RUN go test ./...
RUN swag init
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-extldflags=-static" -o bin/app

############################
# STEP 2 build a small image
############################
FROM alpine

WORKDIR /root/
COPY --from=builder /go/src/qonto-service/test/qonto_accounts.sqlite .
COPY --from=builder /go/src/qonto-service/bin/app .
CMD ["./app"]


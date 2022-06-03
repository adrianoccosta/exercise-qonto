############################
# STEP 1 build executable binary
############################
FROM golang:1.18.3 AS builder

WORKDIR /go/src/qonto-service

COPY . .

RUN rm -rf ./bin
RUN go clean -i ./...
RUN go mod vendor
RUN CGO_ENABLED=1 GOOS=linux go build -o bin/app

############################
# STEP 2 build a small image
############################
FROM alpine

WORKDIR /root/
COPY --from=builder /go/src/qonto-service/interview-test-backend-assets/qonto_accounts.sqlite .
COPY --from=builder /go/src/qonto-service/bin/app .
CMD ["./app"]


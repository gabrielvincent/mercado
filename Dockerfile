ARG GO_VERSION=1
FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /app
COPY . .

RUN go mod download && go mod verify
RUN go build -v -o run-app .

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/run-app /app/run-app
COPY --from=builder /app/public /app/public

WORKDIR /app

RUN chmod +x /app/run-app

CMD ["./run-app"]

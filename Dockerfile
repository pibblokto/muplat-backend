FROM golang:1.24-alpine3.21 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o /app/muplat ./cmd/main.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/muplat /app/muplat

CMD ["./muplat"]

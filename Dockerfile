FROM golang:1.16.0-alpine3.13 as builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build cmd/main.go

FROM gcr.io/distroless/base:latest

WORKDIR /app

COPY --from=builder /go/src/app/main /app/main

CMD ["./main"]

FROM golang:1.25.3-trixie AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./cmd/main .

FROM alpine:3.22.1

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
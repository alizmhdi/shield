FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o shield ./cmd

FROM alpine:latest

COPY --from=builder /app/shield /bin/

RUN chmod +x /bin/shield

ENTRYPOINT ["/bin/shield"]
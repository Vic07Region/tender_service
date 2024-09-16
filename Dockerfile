FROM golang:1.23.0-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o main ./cmd/main.go

FROM ubuntu:latest AS runner

COPY --from=builder /app/main .
COPY --from=builder /app/root.crt .

CMD ["/main"]

EXPOSE 8080
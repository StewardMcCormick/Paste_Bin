FROM golang:1.25.3-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/app ./cmd/app/main.go

FROM alpine
WORKDIR /app

COPY --from=builder /bin/app /app-bin

COPY config.yaml ./

ENTRYPOINT ["/app-bin"]

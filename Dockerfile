FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build .

FROM docker:28-cli
COPY --from=builder /app/image-availability-exporter /usr/bin/image-availability-exporter
ENTRYPOINT ["/usr/bin/image-availability-exporter"]

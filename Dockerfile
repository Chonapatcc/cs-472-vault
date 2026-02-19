FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies needed for downloading modules and certs
RUN apk add --no-cache ca-certificates git

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy sources
COPY . .

# Build statically for a small final image
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /go-app ./

FROM scratch

# Copy binary and CA certs for TLS
COPY --from=builder /go-app /go-app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8008

USER 65532:65532 

ENTRYPOINT ["/go-app"]

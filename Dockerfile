FROM golang:1.19.3-alpine3.16 as builder

WORKDIR /workdir

# Download deps
COPY . .
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "bin/app"

# Deploy
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workdir/bin/app .
USER nonroot:nonroot

EXPOSE 8088

ENTRYPOINT ["./app"]
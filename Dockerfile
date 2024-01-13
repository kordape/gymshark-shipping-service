FROM golang:1.21 as builder
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service cmd/main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/service /service
EXPOSE 80
# Run the web service on container startup.
ENTRYPOINT ["/service"]
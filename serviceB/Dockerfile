FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -C "cmd/server" -ldflags="-w -s" -o server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/cmd/server .
CMD ["./server"]

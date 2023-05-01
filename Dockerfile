FROM golang:alpine AS builder

WORKDIR /app
COPY . .
RUN go build -v -o build/server server/main.go


FROM alpine
COPY --from=builder /app/build/server /app/server

ENTRYPOINT ["/app/server"]
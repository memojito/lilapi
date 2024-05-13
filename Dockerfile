FROM golang:alpine AS builder

WORKDIR /app
COPY . .
RUN go build -v -o build/server server/main.go


FROM surnet/alpine-wkhtmltopdf:3.19.1-0.12.6-small
COPY --from=builder /app/build/server /app/server

ENTRYPOINT ["/app/server"]
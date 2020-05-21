FROM golang:1.14.3-alpine3.11 AS builder
RUN apk add --no-cache upx=~3.95-r2
WORKDIR /app
COPY . .
RUN go build -ldflags="-s -w" && upx rest-geoip

FROM alpine:3.11.6
COPY --from=builder /app/rest-geoip /usr/bin/
ENV GIN_MODE=release
CMD ["/usr/bin/rest-geoip"]

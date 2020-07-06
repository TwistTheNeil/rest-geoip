# Build app
FROM golang:1.14.3-alpine3.11 AS builder
RUN apk add --no-cache upx=~3.95-r2
WORKDIR /app
COPY . .
RUN sh init_deps.sh
RUN go get github.com/markbates/pkger/cmd/pkger && $GOPATH/bin/pkger
RUN go get github.com/securego/gosec/cmd/gosec && $GOPATH/bin/gosec ./...
RUN go build -v -ldflags="-s -w" && upx rest-geoip

# Main docker image
FROM alpine:3.11.6
COPY --from=builder /app/rest-geoip /usr/bin/
ENV GIN_MODE=release
CMD ["/usr/bin/rest-geoip"]

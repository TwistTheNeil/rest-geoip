# Download and verify the GeoLite2 database
FROM alpine:3.11.6 AS maxmind-downloader
WORKDIR /opt

ARG MAXMIND_LICENSE

RUN apk --update add tar
RUN wget "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=${MAXMIND_LICENSE}&suffix=tar.gz" -O GeoLite2-City.tar.gz && \
      wget "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=${MAXMIND_LICENSE}&suffix=tar.gz.md5" -O GeoLite2-City.tar.gz.md5 && \
      EXPECTED_MD5=$(cat GeoLite2-City.tar.gz.md5) && \
      echo "$EXPECTED_MD5  GeoLite2-City.tar.gz" > expected.md5 && \
      md5sum -c expected.md5 && \
      tar -xvf GeoLite2-City.tar.gz --no-anchored 'GeoLite2-City.mmdb' --strip-components=1

# Build app
FROM golang:1.14.3-alpine3.11 AS builder
RUN apk add --no-cache upx=~3.95-r2
WORKDIR /app
COPY . .
RUN go get github.com/securego/gosec/cmd/gosec
RUN $GOPATH/bin/gosec ./...
RUN go build -ldflags="-s -w" && upx rest-geoip

# Main docker image
FROM alpine:3.11.6
COPY --from=maxmind-downloader /opt/GeoLite2-City.mmdb /opt/
COPY --from=builder /app/rest-geoip /usr/bin/
ENV GIN_MODE=release
CMD ["/usr/bin/rest-geoip"]

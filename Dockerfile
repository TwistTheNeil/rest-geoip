FROM node:24.16.0-alpine3.24 AS node

FROM golang:1.26.4-alpine3.24 AS golang

FROM node AS frontend-builder
WORKDIR /app
RUN apk add pnpm
COPY frontend /app
RUN rm -rf /app/node_modules
RUN pnpm install --frozen-lockfile

# Build spa
FROM frontend-builder AS spa-builder
RUN npx vite build --outDir /app/dist

# Build app
FROM golang AS builder
RUN apk add --no-cache upx=5.2.0-r0
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN rm -rf /app/internal/router/dist
COPY --from=spa-builder /app/dist /app/internal/router/dist
RUN go build -v -ldflags="-s"
RUN upx /app/rest-geoip

# dev docker image
FROM golang AS dev
RUN go install github.com/air-verse/air@v1.65.3
EXPOSE 1323
WORKDIR /app

# Main docker image
FROM alpine:3.24
COPY --from=builder /app/rest-geoip /usr/bin/
ENV RELEASE_MODE=true
EXPOSE 1323
CMD ["/usr/bin/rest-geoip"]

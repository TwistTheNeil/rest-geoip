FROM node:20.11.1-alpine3.19 AS frontend-builder
WORKDIR /app
RUN npm install -g pnpm@8.10.0
COPY frontend /app
RUN rm -rf /app/node_modules
RUN pnpm install --frozen-lockfile

# Build spa
FROM frontend-builder AS spa-builder
RUN npx vite build --outDir /app/dist

# Build app
FROM golang:1.21.3-alpine3.18 AS builder
RUN apk add --no-cache upx=4.0.2-r0
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
FROM golang:1.21.3-alpine3.18 AS dev
RUN go install github.com/cosmtrek/air@latest
EXPOSE 1323
WORKDIR /app

# Main docker image
FROM alpine:3.18.4
COPY --from=builder /app/rest-geoip /usr/bin/
ENV RELEASE_MODE=true
EXPOSE 1323
CMD ["/usr/bin/rest-geoip"]

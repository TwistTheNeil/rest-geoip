# Build spa
FROM node:18.18.2-alpine3.18 AS spa-builder
WORKDIR /app
COPY frontend /app
RUN npm install
RUN npx vite build --outDir /app/dist

# Build app
FROM golang:1.21.3-alpine3.17 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
COPY --from=spa-builder /app/dist /app/internal/router/dist
RUN go build -v 

# Main docker image
FROM alpine:3.18.4
COPY --from=builder /app/rest-geoip /usr/bin/
ENV RELEASE_MODE=true
CMD ["/usr/bin/rest-geoip"]
